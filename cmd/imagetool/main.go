package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/Cloud-Foundations/Dominator/lib/constants"
	"github.com/Cloud-Foundations/Dominator/lib/filesystem"
	"github.com/Cloud-Foundations/Dominator/lib/filter"
	"github.com/Cloud-Foundations/Dominator/lib/flags/commands"
	"github.com/Cloud-Foundations/Dominator/lib/flags/loadflags"
	"github.com/Cloud-Foundations/Dominator/lib/flagutil"
	"github.com/Cloud-Foundations/Dominator/lib/log"
	"github.com/Cloud-Foundations/Dominator/lib/log/cmdlogger"
	"github.com/Cloud-Foundations/Dominator/lib/mbr"
	"github.com/Cloud-Foundations/Dominator/lib/net/rrdialer"
	objectclient "github.com/Cloud-Foundations/Dominator/lib/objectserver/client"
	"github.com/Cloud-Foundations/Dominator/lib/srpc"
	"github.com/Cloud-Foundations/Dominator/lib/srpc/setupclient"
	"github.com/Cloud-Foundations/Dominator/lib/tags"
)

var (
	allocateBlocks = flag.Bool("allocateBlocks", false,
		"If true, allocate blocks when making raw image")
	buildCommitId = flag.String("buildCommitId", "",
		"build Commit Id to match when finding latest image")
	buildLog = flag.String("buildLog", "",
		"Filename or URL containing build log")
	compress      = flag.Bool("compress", false, "If true, compress tar output")
	computedFiles = flag.String("computedFiles", "",
		"Name of file containing computed files list")
	computedFilesRoot = flag.String("computedFilesRoot", "",
		"Name of directory tree containing computed files to replace on unpack")
	copyMtimesFrom = flag.String("copyMtimesFrom", "",
		"Name of image to copy mtimes for otherwise unchanged files/devices")
	debug = flag.Bool("debug", false,
		"If true, show debugging output")
	deleteFilter = flag.String("deleteFilter", "",
		"Name of delete filter file for addi, adds and diff subcommands")
	expiresIn = flag.Duration("expiresIn", 0,
		"How long before the image expires (auto deletes). Default: never")
	filterFile = flag.String("filterFile", "",
		"Filter file to apply when adding, diffing or showing images")
	fleetManagerHostname = flag.String("fleetManagerHostname", "",
		"Hostname of Fleet Manager (to find VM to scan)")
	fleetManagerPortNum = flag.Uint("fleetManagerPortNum",
		constants.FleetManagerPortNumber, "Port number of Fleet Manager")
	diffArgs           flagutil.StringList
	hypervisorHostname = flag.String("hypervisorHostname", "",
		"Hostname of hypervisor (for VM to scan)")
	hypervisorPortNum = flag.Uint("hypervisorPortNum",
		constants.HypervisorPortNumber, "Port number of hypervisor")
	ignoreExpiring = flag.Bool("ignoreExpiring", false,
		"If true, ignore expiring images when finding images")
	ignoreFilters = flag.Bool("ignoreFilters", false,
		"If true, ignore filter(s) when diffing/patching")
	imageServerHostname = flag.String("imageServerHostname", "localhost",
		"Hostname of image server")
	imageServerPortNum = flag.Uint("imageServerPortNum",
		constants.ImageServerPortNumber,
		"Port number of image server")
	makeBootable = flag.Bool("makeBootable", true,
		"If true, make raw image bootable by installing GRUB")
	masterImageServerHostname = flag.String("masterImageServerHostname", "",
		"Hostname of master image server (if different)")
	mdbServerHostname = flag.String("mdbServerHostname", "localhost",
		"Hostname of MDB server")
	mdbServerPortNum = flag.Uint("mdbServerPortNum",
		constants.SimpleMdbServerPortNumber,
		"Port number of MDB server")
	minFreeBytes      flagutil.Size = 4 << 20
	objectAddInterval               = flag.Duration("objectAddInterval", 0,
		"Interval between object uploads (for debugging)")
	overlayDirectory = flag.String("overlayDirectory", "",
		"Directory tree of files to overlay on top of the image when making raw image")
	releaseNotes = flag.String("releaseNotes", "",
		"Filename or URL containing release notes")
	requiredPaths = flagutil.StringToRuneMap(constants.RequiredPaths)
	rootLabel     = flag.String("rootLabel", "",
		"Label to write for root file-system when making raw image")
	roundupPower = flag.Uint64("roundupPower", 24,
		"power of 2 to round up raw image size")
	runTriggers = flag.Bool("runTriggers", false,
		"If true, run image triggers when patching /")
	scanExcludeList flagutil.StringList = constants.ScanExcludeList
	skipFields                          = flag.String("skipFields", "",
		"Fields to skip when showing or diffing images")
	tableType   mbr.TableType = mbr.TABLE_TYPE_MSDOS
	tagsToMatch tags.MatchTags
	timeout     = flag.Duration("timeout", 0,
		"Timeout for get and wait subcommands")

	logger            log.DebugLogger
	minimumExpiration = 15 * time.Minute
	rrDialer          *rrdialer.Dialer
)

func init() {
	flag.Var(&diffArgs, "diffArgs",
		"Comma separated list of optional arguments to pass to diffing tool")
	flag.Var(&minFreeBytes, "minFreeBytes",
		"minimum number of free bytes in raw image")
	flag.Var(&requiredPaths, "requiredPaths",
		"Comma separated list of required path:type entries")
	flag.Var(&scanExcludeList, "scanExcludeList",
		"Comma separated list of patterns to exclude from scanning")
	flag.Var(&tableType, "tableType", "partition table type for make-raw-image")
	flag.Var(&tagsToMatch, "tagsToMatch", "Tags to match when finding/listing")
}

func printUsage() {
	w := flag.CommandLine.Output()
	fmt.Fprintln(w,
		"Usage: imagetool [flags...] add|check|delete|list [args...]")
	fmt.Fprintln(w, "Common flags:")
	flag.PrintDefaults()
	fmt.Fprintln(w, "Commands:")
	commands.PrintCommands(w, subcommands)
	fmt.Fprintln(w, "Images can be specified as name:type. Supported types:")
	fmt.Fprintln(w, "  d: name of directory tree to scan")
	fmt.Fprintln(w, "  f: name of file containing a FileSystem")
	fmt.Fprintln(w, "  i: name of an image on the imageserver")
	fmt.Fprintln(w, "  I: name of an image stream on the imageserver (latest)")
	fmt.Fprintln(w, "  l: name of file containing an Image")
	fmt.Fprintln(w, "  s: name of sub to poll")
	fmt.Fprintln(w, "  v: hostname/IP of SmallStack VM to scan")
	fmt.Fprintln(w, "SkipFields:")
	fmt.Fprintln(w, "  m: mode")
	fmt.Fprintln(w, "  l: number of hardlinks")
	fmt.Fprintln(w, "  u: UID")
	fmt.Fprintln(w, "  g: GID")
	fmt.Fprintln(w, "  s: size/Rdev")
	fmt.Fprintln(w, "  t: time of last modification")
	fmt.Fprintln(w, "  n: name")
	fmt.Fprintln(w, "  d: data (hash or symlink target)")
}

var subcommands = []commands.Command{
	{"add", "                    name imagefile filterfile triggerfile", 4, 4,
		addImagefileSubcommand},
	{"addi", "                   name imagename filterfile triggerfile", 4, 4,
		addImageimageSubcommand},
	{"addrep", "                 name baseimage layerimage...", 3, -1,
		addReplaceImageSubcommand},
	{"adds", "                   name subname filterfile triggerfile", 4, 4,
		addImagesubSubcommand},
	{"analyse-file-system", "    directory", 1, 1, analyseFileSystemSubcommand},
	{"bulk-addrep", "            layerimage...", 1, -1,
		bulkAddReplaceImagesSubcommand},
	{"change-image-expiration", "name", 1, 1, changeImageExpirationSubcommand},
	{"check", "                  name", 1, 1, checkImageSubcommand},
	{"check-directory", "        dirname", 1, 1, checkDirectorySubcommand},
	{"chown", "                  dirname ownerGroup", 2, 2,
		chownDirectorySubcommand},
	{"copy", "                   name oldimagename", 2, 2, copyImageSubcommand},
	{"copy-filtered-files", "    name srcdir destdir", 3, 3,
		copyFilteredFilesSubcommand},
	{"delete", "                 name", 1, 1, deleteImageSubcommand},
	{"delunrefobj", "            percentage bytes", 2, 2,
		deleteUnreferencedObjectsSubcommand},
	{"diff", "                   tool left right", 3, 3, diffSubcommand},
	{"diff-build-logs", "        tool left right", 3, 3,
		diffBuildLogsInImagesSubcommand},
	{"diff-files", "             tool left right filename", 4, 4,
		diffFileInImagesSubcommand},
	{"diff-filters", "           tool left right", 3, 3,
		diffFilterInImagesSubcommand},
	{"diff-package-lists", "     tool left right", 3, 3,
		diffImagePackageListsSubcommand},
	{"diff-triggers", "          tool left right", 3, 3,
		diffTriggersInImagesSubcommand},
	{"estimate-usage", "         name", 1, 1, estimateImageUsageSubcommand},
	{"find-latest-image", "      directory", 1, 1, findLatestImageSubcommand},
	{"get", "                    name directory", 2, 2, getImageSubcommand},
	{"get-archive-data", "       name outfile", 2, 2,
		getImageArchiveDataSubcommand},
	{"get-build-log", "          name [outfile]", 1, 2,
		getImageBuildLogSubcommand},
	{"get-file-in-image", "      name imageFile [outfile]", 2, 3,
		getFileInImageSubcommand},
	{"get-image-expiration", "   name", 1, 1, getImageExpirationSubcommand},
	{"get-image-updates", "", 0, 0, getImageUpdatesSubcommand},
	{"get-package-list", "       name [outfile]", 1, 2,
		getImagePackageListSubcommand},
	{"get-replication-master", "", 0, 0, getReplicationMasterSubcommand},
	{"list", "", 0, 0, listImagesSubcommand},
	{"list-mdb", "", 0, 0, listMdbImagesSubcommand},
	{"list-not-in-mdb", "", 0, 0, listImagesNotInMdbSubcommand},
	{"listdirs", "", 0, 0, listDirectoriesSubcommand},
	{"listunrefobj", "", 0, 0, listUnreferencedObjectsSubcommand},
	{"make-raw-image", "         name rawfile", 2, 2, makeRawImageSubcommand},
	{"match-triggers", "         name triggers-file", 2, 2,
		matchTriggersSubcommand},
	{"merge-filters", "          filter-file...", 1, -1, mergeFiltersSubcommand},
	{"merge-triggers", "         triggers-file...", 1, -1,
		mergeTriggersSubcommand},
	{"mkdir", "                  name", 1, 1, makeDirectorySubcommand},
	{"patch-directory", "        name directory", 2, 2,
		patchDirectorySubcommand},
	{"restore-from-file", "      filename", 1, 1, restoreImageSubcommand},
	{"save-to-file", "           name [outfile]", 1, 2, saveImageSubcommand},
	{"scan-filtered-files", "    name directory", 2, 2,
		scanFilteredFilesSubcommand},
	{"show", "                   name", 1, 1, showImageSubcommand},
	{"show-bad-computed-files", "", 0, 0, showBadComputedFilesSubcommand},
	{"show-bad-image-subs", "", 0, 0, showBadImageSubsSubcommand},
	{"show-computed-file-subs", "filename source", 2, 2,
		showComputedFileSubsSubcommand},
	{"show-filter", "            name", 1, 1, showImageFilterSubcommand},
	{"show-inode", "             name inodePath", 2, 2,
		showImageInodeSubcommand},
	{"show-metadata", "          name", 1, 1, showImageMetadataSubcommand},
	{"show-triggers", "          name", 1, 1, showImageTriggersSubcommand},
	{"showunrefobj", "", 0, 0, showUnreferencedObjectsSubcommand},
	{"tar", "                    name [file]", 1, 2, tarImageSubcommand},
	{"test-download-speed", "    name", 1, 1, testDownloadSpeedSubcommand},
	{"trace-inode-history", "    name inodePath", 2, 2,
		traceInodeHistorySubcommand},
	{"wait", "                   name", 1, 1, waitImageSubcommand},
}

var (
	imageSrpcClient       *srpc.Client
	masterImageSrpcClient *srpc.Client
	theObjectClient       *objectclient.ObjectClient
	theMasterObjectClient *objectclient.ObjectClient

	listSelector filesystem.ListSelector
)

func getClients() (*srpc.Client, *objectclient.ObjectClient) {
	getPointedClients(*imageServerHostname, &imageSrpcClient, &theObjectClient)
	return imageSrpcClient, theObjectClient
}

func getMasterClients() (*srpc.Client, *objectclient.ObjectClient) {
	if *masterImageServerHostname == "" {
		return getClients()
	}
	getPointedClients(*masterImageServerHostname,
		&masterImageSrpcClient, &theMasterObjectClient)
	return masterImageSrpcClient, theMasterObjectClient
}

func getPointedClients(hostname string, iClient **srpc.Client,
	oClient **objectclient.ObjectClient) {
	if *iClient == nil {
		var err error
		clientName := fmt.Sprintf("%s:%d", hostname, *imageServerPortNum)
		*iClient, err = srpc.DialHTTPWithDialer("tcp", clientName, rrDialer)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error dialing: %s: %s\n", clientName, err)
			os.Exit(1)
		}
		*oClient = objectclient.AttachObjectClient(imageSrpcClient)
	}
}

func makeListSelector(arg string) filesystem.ListSelector {
	var mask filesystem.ListSelector = filesystem.ListSelectAll
	for _, char := range arg {
		switch char {
		case 'm':
			mask |= filesystem.ListSelectSkipMode
		case 'l':
			mask |= filesystem.ListSelectSkipNumLinks
		case 'u':
			mask |= filesystem.ListSelectSkipUid
		case 'g':
			mask |= filesystem.ListSelectSkipGid
		case 's':
			mask |= filesystem.ListSelectSkipSizeDevnum
		case 't':
			mask |= filesystem.ListSelectSkipMtime
		case 'n':
			mask |= filesystem.ListSelectSkipName
		case 'd':
			mask |= filesystem.ListSelectSkipData
		}
	}
	return mask
}

var listFilter *filter.Filter

func doMain() int {
	if err := loadflags.LoadForCli("imagetool"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	flag.Usage = printUsage
	flag.Parse()
	if flag.NArg() < 1 {
		printUsage()
		return 2
	}
	logger = cmdlogger.New()
	srpc.SetDefaultLogger(logger)
	if *expiresIn > 0 && *expiresIn < minimumExpiration {
		fmt.Fprintf(os.Stderr, "Minimum expiration: %s\n", minimumExpiration)
		return 2
	}
	listSelector = makeListSelector(*skipFields)
	var err error
	if *filterFile != "" {
		listFilter, err = filter.Load(*filterFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 2
		}
	}
	if err := setupclient.SetupTls(true); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	rrDialer, err = rrdialer.New(&net.Dialer{Timeout: time.Second * 10}, "",
		logger)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	defer rrDialer.WaitForBackgroundResults(time.Second)
	return commands.RunCommands(subcommands, printUsage, logger)
}

func main() {
	os.Exit(doMain())
}
