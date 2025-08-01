package hypervisor

import (
	"net"
	"time"

	"github.com/Cloud-Foundations/Dominator/lib/filesystem"
	"github.com/Cloud-Foundations/Dominator/lib/filter"
	"github.com/Cloud-Foundations/Dominator/lib/tags"
)

const (
	ConsoleNone  = 0
	ConsoleDummy = 1
	ConsoleVNC   = 2

	FirmwareDefault = 0
	FirmwareBIOS    = 1
	FirmwareUEFI    = 2

	MachineTypeGenericPC = 0
	MachineTypeQ35       = 1

	StateStarting      = 0
	StateRunning       = 1
	StateFailedToStart = 2
	StateStopping      = 3
	StateStopped       = 4
	StateDestroying    = 5
	StateMigrating     = 6
	StateExporting     = 7
	StateCrashed       = 8
	StateDebugging     = 9

	VolumeFormatRaw   = 0
	VolumeFormatQCOW2 = 1

	VolumeInterfaceVirtIO = 0
	VolumeInterfaceIDE    = 1
	VolumeInterfaceNVMe   = 2

	VolumeTypePersistent = 0
	VolumeTypeMemory     = 1

	WatchdogActionNone     = 0
	WatchdogActionReset    = 1
	WatchdogActionShutdown = 2
	WatchdogActionPowerOff = 3

	WatchdogModelNone     = 0
	WatchdogModelIb700    = 1
	WatchdogModelI6300esb = 2
)

type AcknowledgeVmRequest struct {
	IpAddress net.IP
}

type AcknowledgeVmResponse struct {
	Error string
}

type Address struct {
	IpAddress  net.IP `json:",omitempty"`
	MacAddress string
}

type AddressList []Address

type AddVmVolumesRequest struct {
	IpAddress   net.IP
	VolumeSizes []uint64
}

type AddVmVolumesResponse struct {
	Error string
}

type BecomePrimaryVmOwnerRequest struct {
	IpAddress net.IP
}

type BecomePrimaryVmOwnerResponse struct {
	Error string
}

type ChangeAddressPoolRequest struct {
	AddressesToAdd       []Address       // Will be added to free pool.
	AddressesToRemove    []Address       // Will be removed from free pool.
	MaximumFreeAddresses map[string]uint // Key: subnet ID.
}

type ChangeAddressPoolResponse struct {
	Error string
}

type ChangeOwnersRequest struct {
	OwnerGroups []string `json:",omitempty"`
	OwnerUsers  []string `json:",omitempty"`
}

type ChangeOwnersResponse struct {
	Error string
}

type ChangeVmConsoleTypeRequest struct {
	ConsoleType ConsoleType
	IpAddress   net.IP
}

type ChangeVmConsoleTypeResponse struct {
	Error string
}

type ChangeVmCpuPriorityRequest struct {
	CpuPriority int
	IpAddress   net.IP
}

type ChangeVmCpuPriorityResponse struct {
	Error string
}

type ChangeVmDestroyProtectionRequest struct {
	DestroyProtection bool
	IpAddress         net.IP
}

type ChangeVmDestroyProtectionResponse struct {
	Error string
}

type ChangeVmMachineTypeRequest struct {
	MachineType MachineType
	IpAddress   net.IP
}

type ChangeVmMachineTypeResponse struct {
	Error string
}

type ChangeVmOwnerGroupsRequest struct {
	IpAddress   net.IP
	OwnerGroups []string
}

type ChangeVmOwnerGroupsResponse struct {
	Error string
}

type ChangeVmOwnerUsersRequest struct {
	IpAddress  net.IP
	OwnerUsers []string
}

type ChangeVmOwnerUsersResponse struct {
	Error string
}

type ChangeVmSizeRequest struct {
	IpAddress   net.IP
	MemoryInMiB uint64
	MilliCPUs   uint
	VirtualCPUs uint
}

type ChangeVmSubnetRequest struct {
	IpAddress net.IP
	SubnetId  string
}

type ChangeVmSubnetResponse struct {
	Error           string
	NewIpAddress    net.IP
	OldIdentityName string
}

type ChangeVmSizeResponse struct {
	Error string
}

type ChangeVmTagsRequest struct {
	IpAddress net.IP
	Tags      tags.Tags
}

type ChangeVmTagsResponse struct {
	Error string
}

type ChangeVmVolumeInterfacesRequest struct {
	Interfaces []VolumeInterface
	IpAddress  net.IP
}

type ChangeVmVolumeInterfacesResponse struct {
	Error string
}

type ChangeVmVolumeSizeRequest struct {
	IpAddress   net.IP
	VolumeIndex uint
	VolumeSize  uint64
}

type ChangeVmVolumeSizeResponse struct {
	Error string
}

type CommitImportedVmRequest struct {
	IpAddress net.IP
}

type CommitImportedVmResponse struct {
	Error string
}

// The ConnectToVmConsole RPC is fully streamed. After the request/response,
// the connection/client is hijacked and each side of the connection will send
// a stream of bytes.
type ConnectToVmConsoleRequest struct {
	IpAddress net.IP
}

type ConnectToVmConsoleResponse struct {
	Error string
}

// The ConnectToVmManger RPC is fully streamed. After the request/response,
// the connection/client is hijacked and each side of the connection will send
// a stream of bytes.
type ConnectToVmManagerRequest struct {
	IpAddress net.IP
}

type ConnectToVmManagerResponse struct {
	Error string
}

// The ConnectToVmSerialPort RPC is fully streamed. After the request/response,
// the connection/client is hijacked and each side of the connection will send
// a stream of bytes.
type ConnectToVmSerialPortRequest struct {
	IpAddress  net.IP
	PortNumber uint
}

type ConnectToVmSerialPortResponse struct {
	Error string
}

type ConsoleType uint

type CopyVmRequest struct {
	AccessToken      []byte
	IpAddress        net.IP
	SkipMemoryCheck  bool
	SourceHypervisor string
	VmInfo
}

type CopyVmResponse struct { // Multiple responses are sent.
	Error           string
	Final           bool // If true, this is the final response.
	IpAddress       net.IP
	ProgressMessage string
}

type CreateVmRequest struct {
	DhcpTimeout          time.Duration // <0: no DHCP; 0: no wait; >0 DHPC wait.
	DoNotStart           bool
	EnableNetboot        bool
	IdentityCertificate  []byte // PEM encoded.
	IdentityKey          []byte // PEM encoded.
	ImageDataSize        uint64
	ImageTimeout         time.Duration
	MinimumFreeBytes     uint64
	OverlayDirectories   []string
	OverlayFiles         map[string][]byte
	RoundupPower         uint64
	SecondaryVolumes     []Volume
	SecondaryVolumesData bool // Exclusive of SecondaryVolumesInit.
	SecondaryVolumesInit []VolumeInitialisationInfo
	SkipBootloader       bool
	SkipMemoryCheck      bool
	StorageIndices       []uint
	UserDataSize         uint64
	VmInfo
} // The following data are streamed afterwards in the following order:
//     RAW image data (length=ImageDataSize)
//     user data (length=UserDataSize)
//     secondary volumes (if SecondaryVolumesData is true)

type CreateVmResponse struct { // Multiple responses are sent.
	DhcpTimedOut    bool
	Final           bool // If true, this is the final response.
	IpAddress       net.IP
	ProgressMessage string
	Error           string
}

type DebugVmImageRequest struct {
	DhcpTimeout      time.Duration // <0: no DHCP; 0: no wait; >0 DHPC wait.
	ImageDataSize    uint64
	ImageName        string
	ImageTimeout     time.Duration
	ImageURL         string
	IpAddress        net.IP
	MinimumFreeBytes uint64
	OverlayFiles     map[string][]byte
	RoundupPower     uint64
} // The following data are streamed afterwards in the following order:
//     RAW image data (length=ImageDataSize)

type DebugVmImageResponse struct { // Multiple responses are sent.
	DhcpTimedOut    bool
	Final           bool // If true, this is the final response.
	ProgressMessage string
	Error           string
}

type DeleteVmVolumeRequest struct {
	AccessToken []byte
	IpAddress   net.IP
	VolumeIndex uint
}

type DeleteVmVolumeResponse struct {
	Error string
}

type DestroyVmRequest struct {
	AccessToken []byte
	IpAddress   net.IP
}

type DestroyVmResponse struct {
	Error string
}

type DiscardVmAccessTokenRequest struct {
	AccessToken []byte
	IpAddress   net.IP
}

type DiscardVmAccessTokenResponse struct {
	Error string
}

type DiscardVmOldImageRequest struct {
	IpAddress net.IP
}

type DiscardVmOldImageResponse struct {
	Error string
}

type DiscardVmOldUserDataRequest struct {
	IpAddress net.IP
}

type DiscardVmOldUserDataResponse struct {
	Error string
}

type DiscardVmSnapshotRequest struct {
	IpAddress net.IP
	Name      string
}

type DiscardVmSnapshotResponse struct {
	Error string
}

type ExportLocalVmInfo struct {
	Bridges []string
	LocalVmInfo
}

type ExportLocalVmRequest struct {
	IpAddress          net.IP
	VerificationCookie []byte `json:",omitempty"`
}

type ExportLocalVmResponse struct {
	Error  string
	VmInfo ExportLocalVmInfo
}

type FirmwareType uint

type GetCapacityRequest struct{}

type GetCapacityResponse struct {
	MemoryInMiB      uint64 `json:",omitempty"`
	NumCPUs          uint   `json:",omitempty"`
	TotalVolumeBytes uint64 `json:",omitempty"`
}

type GetIdentityProviderRequest struct{}

type GetIdentityProviderResponse struct {
	Error   string
	BaseUrl string
}

type GetPublicKeyRequest struct{}

type GetPublicKeyResponse struct {
	Error  string
	KeyPEM []byte
}

type GetRootCookiePathRequest struct{}

type GetRootCookiePathResponse struct {
	Error string
	Path  string
}

// The GetUpdates() RPC is fully streamed.
// The client may or may not send GetUpdatesRequest messages to the server.
// The server sends a stream of Update messages.

type GetUpdatesRequest struct {
	RegisterExternalLeasesRequest *RegisterExternalLeasesRequest
}

type Update struct {
	HaveAddressPool  bool               `json:",omitempty"`
	AddressPool      []Address          `json:",omitempty"` // Used & free.
	HaveDisabled     bool               `json:",omitempty"`
	Disabled         bool               `json:",omitempty"`
	MemoryInMiB      *uint64            `json:",omitempty"`
	NumCPUs          *uint              `json:",omitempty"`
	NumFreeAddresses map[string]uint    `json:",omitempty"` // Key: subnet ID.
	HealthStatus     string             `json:",omitempty"`
	HaveSerialNumber bool               `json:",omitempty"`
	SerialNumber     string             `json:",omitempty"`
	HaveSubnets      bool               `json:",omitempty"`
	Subnets          []Subnet           `json:",omitempty"`
	TotalVolumeBytes *uint64            `json:",omitempty"`
	HaveVMs          bool               `json:",omitempty"`
	VMs              map[string]*VmInfo `json:",omitempty"` // Key: IP address.
}

type GetVmAccessTokenRequest struct {
	IpAddress net.IP
	Lifetime  time.Duration
}

type GetVmAccessTokenResponse struct {
	Token []byte `json:",omitempty"`
	Error string
}

type GetVmInfoRequest struct {
	IpAddress net.IP
}

type GetVmInfoResponse struct {
	VmInfo VmInfo
	Error  string
}

type GetVmInfosRequest struct {
	IgnoreStateMask uint64
	OwnerGroups     []string
	OwnerUsers      []string
	VmTagsToMatch   tags.MatchTags // Empty: match all tags.
}

type GetVmInfosResponse struct {
	Error   string
	VmInfos []VmInfo
}

type GetVmLastPatchLogRequest struct {
	IpAddress net.IP
}

type GetVmLastPatchLogResponse struct {
	Error     string
	Length    uint64
	PatchTime time.Time
} // Data (length=Length) are streamed afterwards.

type GetVmUserDataRequest struct {
	AccessToken []byte
	IpAddress   net.IP
}

type GetVmUserDataResponse struct {
	Error  string
	Length uint64
} // Data (length=Length) are streamed afterwards.

// The GetVmVolume() RPC is followed by the proto/rsync.GetBlocks message.

type GetVmVolumeRequest struct {
	AccessToken      []byte
	GetExtraFiles    bool
	IgnoreExtraFiles bool
	IpAddress        net.IP
	VolumeIndex      uint
}

type GetVmVolumeResponse struct {
	Error      string
	ExtraFiles map[string][]byte // May contain "kernel", "initrd" and such.
}

type HoldLockRequest struct {
	Timeout   time.Duration
	WriteLock bool
}

type HoldLockResponse struct {
	Error string
}

type HoldVmLockRequest struct {
	IpAddress net.IP
	Timeout   time.Duration
	WriteLock bool
}

type HoldVmLockResponse struct {
	Error string
}

type ListSubnetsRequest struct {
	Sort bool
}

type ListSubnetsResponse struct {
	Error   string
	Subnets []Subnet `json:",omitempty"`
}

type ImportLocalVmRequest struct {
	SkipMemoryCheck    bool
	VerificationCookie []byte `json:",omitempty"`
	VmInfo
	VolumeFilenames []string
}

type ImportLocalVmResponse struct {
	Error string
}

type ListVMsRequest struct {
	IgnoreStateMask uint64
	OwnerGroups     []string
	OwnerUsers      []string
	Sort            bool
	VmTagsToMatch   tags.MatchTags // Empty: match all tags.
}

type ListVMsResponse struct {
	IpAddresses []net.IP
}

type ListVolumeDirectoriesRequest struct{}

type ListVolumeDirectoriesResponse struct {
	Directories []string
	Error       string
}

type LocalVolume struct {
	DirectoryToCleanup string
	Filename           string
}

type LocalVmInfo struct {
	VmInfo
	VolumeLocations []LocalVolume
}

type MachineType uint

type MigrateVmRequest struct {
	AccessToken      []byte
	DhcpTimeout      time.Duration
	IpAddress        net.IP
	SkipMemoryCheck  bool
	SourceHypervisor string
}

type MigrateVmResponse struct { // Multiple responses are sent.
	Error           string
	Final           bool // If true, this is the final response.
	ProgressMessage string
	RequestCommit   bool
}

type MigrateVmResponseResponse struct {
	Commit bool
}

type NetbootMachineRequest struct {
	Address                      Address
	Files                        map[string][]byte
	FilesExpiration              time.Duration
	Hostname                     string
	NumAcknowledgementsToWaitFor uint
	OfferExpiration              time.Duration
	Subnet                       *Subnet
	WaitTimeout                  time.Duration
}

type NetbootMachineResponse struct {
	Error string
}

type PatchVmImageRequest struct {
	ImageName    string
	ImageTimeout time.Duration
	IpAddress    net.IP
	SkipBackup   bool
}

type PatchVmImageResponse struct { // Multiple responses are sent.
	Final           bool // If true, this is the final response.
	ProgressMessage string
	Error           string
}

type PowerOffRequest struct {
	StopVMs bool // true: attempt to stop VMs; false: running VMs block poweroff
}

type PowerOffResponse struct {
	Error string
}

type PrepareVmForMigrationRequest struct {
	AccessToken []byte
	Enable      bool
	IpAddress   net.IP
}

type PrepareVmForMigrationResponse struct {
	Error string
}

type ProbeVmPortRequest struct {
	IpAddress  net.IP
	PortNumber uint
	Timeout    time.Duration
}

type ProbeVmPortResponse struct {
	PortIsOpen bool
	Error      string
}

type RebootVmRequest struct {
	DhcpTimeout time.Duration
	IpAddress   net.IP
}

type RebootVmResponse struct {
	DhcpTimedOut bool
	Error        string
}

type RegisterExternalLeasesRequest struct {
	Addresses AddressList
	Hostnames []string `json:",omitempty"`
}

type RegisterExternalLeasesResponse struct {
	Error string
}

type ReplaceVmCredentialsRequest struct {
	IdentityCertificate []byte // PEM encoded.
	IdentityKey         []byte // PEM encoded.
	IpAddress           net.IP
}

type ReplaceVmCredentialsResponse struct {
	Error string
}

type ReplaceVmIdentityRequest struct {
	IdentityRequestorCertificate []byte // PEM encoded.
	IpAddress                    net.IP
}

type ReplaceVmIdentityResponse struct {
	Error string
}

type ReplaceVmImageRequest struct {
	DhcpTimeout      time.Duration
	ImageDataSize    uint64
	ImageName        string `json:",omitempty"`
	ImageTimeout     time.Duration
	ImageURL         string `json:",omitempty"`
	IpAddress        net.IP
	MinimumFreeBytes uint64
	OverlayFiles     map[string][]byte
	RoundupPower     uint64
	SkipBackup       bool
	SkipBootloader   bool
} // RAW image data (length=ImageDataSize) is streamed afterwards.

type ReplaceVmImageResponse struct { // Multiple responses are sent.
	DhcpTimedOut    bool
	Final           bool // If true, this is the final response.
	ProgressMessage string
	Error           string
}

type ReplaceVmUserDataRequest struct {
	IpAddress net.IP
	Size      uint64
} // User data (length=Size) are streamed afterwards.

type ReplaceVmUserDataResponse struct {
	Error string
}

type RestoreVmFromSnapshotRequest struct {
	IpAddress         net.IP
	ForceIfNotStopped bool
	Name              string
}

type RestoreVmFromSnapshotResponse struct {
	Error string
}

type RestoreVmImageRequest struct {
	IpAddress net.IP
}

type RestoreVmImageResponse struct {
	Error string
}

type RestoreVmUserDataRequest struct {
	IpAddress net.IP
}

type RestoreVmUserDataResponse struct {
	Error string
}

type ReorderVmVolumesRequest struct {
	AccessToken   []byte
	IpAddress     net.IP
	VolumeIndices []uint
}

type ReorderVmVolumesResponse struct {
	Error string
}

type ScanVmRootRequest struct {
	IpAddress net.IP
	Filter    *filter.Filter
}

type ScanVmRootResponse struct {
	Error      string
	FileSystem *filesystem.FileSystem
}

type SetDisabledStateRequest struct {
	Disable bool
}

type SetDisabledStateResponse struct {
	Error string
}

type SnapshotVmRequest struct {
	IpAddress         net.IP
	ForceIfNotStopped bool
	Name              string
	RootOnly          bool
}

type SnapshotVmResponse struct {
	Error string
}

type StartVmRequest struct {
	AccessToken []byte
	DhcpTimeout time.Duration
	IpAddress   net.IP
}

type StartVmResponse struct {
	DhcpTimedOut bool
	Error        string
}

type StopVmRequest struct {
	AccessToken []byte
	IpAddress   net.IP
}

type StopVmResponse struct {
	Error string
}

type State uint

type Subnet struct {
	Id                string
	IpGateway         net.IP
	IpMask            net.IP // net.IPMask can't be JSON {en,de}coded.
	DomainName        string `json:",omitempty"`
	DomainNameServers []net.IP
	DisableMetadata   bool     `json:",omitempty"`
	Manage            bool     `json:",omitempty"`
	VlanId            uint     `json:",omitempty"`
	AllowedGroups     []string `json:",omitempty"`
	AllowedUsers      []string `json:",omitempty"`
	FirstDynamicIP    net.IP   `json:",omitempty"`
	LastDynamicIP     net.IP   `json:",omitempty"`
}

type TraceVmMetadataRequest struct {
	IpAddress net.IP
}

type TraceVmMetadataResponse struct {
	Error string
} // A stream of strings (trace paths) follow.

type UpdateSubnetsRequest struct {
	Add    []Subnet
	Change []Subnet
	Delete []string
}

type UpdateSubnetsResponse struct {
	Error string
}

type VmInfo struct {
	Address             Address
	ChangedStateOn      time.Time    `json:",omitempty"`
	ConsoleType         ConsoleType  `json:",omitempty"`
	CreatedOn           time.Time    `json:",omitempty"`
	CpuPriority         int          `json:",omitempty"`
	DestroyOnPowerdown  bool         `json:",omitempty"`
	DestroyProtection   bool         `json:",omitempty"`
	DisableVirtIO       bool         `json:",omitempty"`
	ExtraKernelOptions  string       `json:",omitempty"`
	FirmwareType        FirmwareType `json:",omitempty"`
	Hostname            string       `json:",omitempty"`
	IdentityExpires     time.Time    `json:",omitempty"`
	IdentityName        string       `json:",omitempty"`
	ImageName           string       `json:",omitempty"`
	ImageURL            string       `json:",omitempty"`
	MachineType         MachineType  `json:",omitempty"`
	MemoryInMiB         uint64
	MilliCPUs           uint
	OwnerGroups         []string `json:",omitempty"`
	OwnerUsers          []string `json:",omitempty"`
	RootFileSystemLabel string   `json:",omitempty"`
	SpreadVolumes       bool     `json:",omitempty"`
	State               State
	SecondaryAddresses  []Address      `json:",omitempty"`
	SecondarySubnetIDs  []string       `json:",omitempty"`
	SubnetId            string         `json:",omitempty"`
	Tags                tags.Tags      `json:",omitempty"`
	Uncommitted         bool           `json:",omitempty"`
	VirtualCPUs         uint           `json:",omitempty"`
	Volumes             []Volume       `json:",omitempty"`
	WatchdogAction      WatchdogAction `json:",omitempty"`
	WatchdogModel       WatchdogModel  `json:",omitempty"`
}

type Volume struct {
	Format    VolumeFormat      `json:",omitempty"`
	Interface VolumeInterface   `json:",omitempty"`
	Size      uint64            `json:",omitempty"`
	Snapshots map[string]uint64 `json:",omitempty"`
	Type      VolumeType        `json:",omitempty"`
}

type VolumeFormat uint

type VolumeInterface uint

type VolumeInitialisationInfo struct {
	BytesPerInode            uint64
	Label                    string
	ReservedBlocksPercentage uint16
}

type VolumeType uint

// The WatchDhcp() RPC is fully streamed.
// The client sends a single WatchDhcpRequest message.
// The server sends a stream of WatchDhcpResponse messages until there is an
// error.

type WatchDhcpRequest struct {
	Interface  string // Default: watch all available interfaces.
	MaxPackets uint64 // Zero means infinite.
}

type WatchDhcpResponse struct {
	Error     string
	Interface string
	Packet    []byte
}

type WatchdogAction uint
type WatchdogModel uint
