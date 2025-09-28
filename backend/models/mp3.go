package models

type AudioMetadata struct {
	SampleRate int
	Channels   int
	BitDepth   int
	Duration   float64
	TotalBytes int
}

type Tag interface {
	GetTitle() string
	GetArtist() string
	GetAlbum() string
	GetYear() string
	GetComment() string
	GetGenre() byte
}

type Frame interface {
	GetHeader() *MP3FrameHeader
	GetData() []byte
	Size() int
}

type ID3v1Tag struct {
	Title   string
	Artist  string
	Album   string
	Year    string
	Comment string
	Genre   byte
}

func (t *ID3v1Tag) GetTitle() string   { return t.Title }
func (t *ID3v1Tag) GetArtist() string  { return t.Artist }
func (t *ID3v1Tag) GetAlbum() string   { return t.Album }
func (t *ID3v1Tag) GetYear() string    { return t.Year }
func (t *ID3v1Tag) GetComment() string { return t.Comment }
func (t *ID3v1Tag) GetGenre() byte     { return t.Genre }

type ID3v2Header struct {
	Version [2]byte
	Flags   byte
	Size    int
}

type MP3FrameHeader struct {
	VersionID     int
	Layer         int
	ProtectionBit bool
	Bitrate       int
	SampleRate    int
	Padding       bool
	ChannelMode   int
	FrameLength   int
}

type MP3Frame struct {
	Header      *MP3FrameHeader
	HeaderBytes []byte
	Data        []byte
}

func (f *MP3Frame) GetHeader() *MP3FrameHeader { return f.Header }
func (f *MP3Frame) GetData() []byte            { return f.Data }
func (f *MP3Frame) Size() int                  { return len(f.HeaderBytes) + len(f.Data) }

type MP3File struct {
	ID3v2     *ID3v2Header
	ID3v2Data []byte
	Frames    []Frame
	ID3v1     Tag
}

func NewMP3File() *MP3File {
	return &MP3File{
		Frames: make([]Frame, 0),
	}
}

func (f *MP3File) AddFrame(frame Frame) {
	f.Frames = append(f.Frames, frame)
}

func (f *MP3File) GetTotalFrames() int {
	return len(f.Frames)
}

func (f *MP3File) HasID3v1() bool {
	return f.ID3v1 != nil
}

func (f *MP3File) HasID3v2() bool {
	return f.ID3v2 != nil
}
