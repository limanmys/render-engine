package connector

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/limanmys/go/postgresql"
	"github.com/pkg/sftp"
	"github.com/tus/tusd/pkg/handler"
)

type SftpStore struct {
	Path   string
	Client *sftp.Client
}

// UseIn sets this store as the core data store in the passed composer and adds
// all possible extension to it.
func (store *SftpStore) UseIn(composer *handler.StoreComposer) {
	composer.UseCore(store)
	composer.UseTerminater(store)
	composer.UseConcater(store)
	composer.UseLengthDeferrer(store)
}

func (store *SftpStore) NewUpload(ctx context.Context, info handler.FileInfo) (handler.Upload, error) {
	fmt.Println("Test")
	err := store.establishConnection(info)
	if err != nil {
		err = fmt.Errorf("Failed to create client: " + err.Error())
		return nil, err
	}
	_ = store.Client.MkdirAll(store.Path)

	id := uid()
	binPath := store.binPath(info.MetaData["filename"])
	info.ID = id
	info.Storage = map[string]string{
		"Type": "sftpstore",
		"Path": binPath,
	}

	file, err := store.Client.OpenFile(binPath, os.O_CREATE|os.O_WRONLY)
	if err != nil {
		err = fmt.Errorf("Err: " + err.Error())
		return nil, err
	}
	err = file.Close()
	if err != nil {
		return nil, err
	}

	upload := &sftpUpload{
		info:     info,
		infoPath: store.infoPath(id),
		binPath:  store.binPath(info.MetaData["filename"]),
		client:   store.Client,
	}

	// writeInfo creates the file by itself if necessary
	err = upload.writeInfo()
	if err != nil {
		return nil, err
	}

	return upload, nil
}

func (store *SftpStore) GetUpload(ctx context.Context, id string) (handler.Upload, error) {
	info := handler.FileInfo{}

	file, err := store.Client.Open(store.infoPath(id))
	if err != nil {
		if os.IsNotExist(err) {
			// Interpret os.ErrNotExist as 404 Not Found
			err = handler.ErrNotFound
		}
		return nil, err
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(file)
	contents := buf.Bytes()
	if err := json.Unmarshal(contents, &info); err != nil {
		return nil, err
	}
	defer file.Close()

	binPath := store.binPath(info.MetaData["filename"])
	infoPath := store.infoPath(id)
	stat, err := store.Client.Stat(binPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Interpret os.ErrNotExist as 404 Not Found
			err = handler.ErrNotFound
		}
		return nil, err
	}

	info.Offset = stat.Size()
	upload := &sftpUpload{
		info:     info,
		infoPath: infoPath,
		binPath:  binPath,
		client:   store.Client,
	}

	return upload, nil
}

func (store *SftpStore) establishConnection(info handler.FileInfo) error {
	fmt.Println("Test2")
	server := postgresql.GetServer(info.MetaData["server_id"])
	val, err := GetConnection(info.MetaData["user_id"], info.MetaData["server_id"], server.IPAddress)
	if err != nil {
		return err
	}

	if !val.CreateFileConnection(info.MetaData["user_id"], info.MetaData["server_id"], server.IPAddress) {
		return fmt.Errorf("could not create file connection")
	}
	store.Client = val.SFTP
	store.Path = info.MetaData["path"]

	return nil
}

func (store SftpStore) AsTerminatableUpload(upload handler.Upload) handler.TerminatableUpload {
	return upload.(*sftpUpload)
}

func (store SftpStore) AsLengthDeclarableUpload(upload handler.Upload) handler.LengthDeclarableUpload {
	return upload.(*sftpUpload)
}

func (store SftpStore) AsConcatableUpload(upload handler.Upload) handler.ConcatableUpload {
	return upload.(*sftpUpload)
}

// uid creates random string and returns it
func uid() string {
	id := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, id)
	if err != nil {
		// This is probably an appropriate way to handle errors from our source
		// for random bits.
		panic(err)
	}
	return hex.EncodeToString(id)
}

// binPath returns the path to the file storing the binary data.
func (store SftpStore) binPath(id string) string {
	return filepath.Join(store.Path, id)
}

// infoPath returns the path to the .info file storing the file's info.
func (store SftpStore) infoPath(id string) string {
	return filepath.Join(store.Path, id+".info")
}

type sftpUpload struct {
	// info stores the current information about the upload
	info handler.FileInfo
	// infoPath is the path to the .info file
	infoPath string
	// binPath is the path to the binary file (which has no extension)
	binPath string
	// client is the sftp client
	client *sftp.Client
}

func (upload *sftpUpload) GetInfo(ctx context.Context) (handler.FileInfo, error) {
	return upload.info, nil
}

func (upload *sftpUpload) WriteChunk(ctx context.Context, offset int64, src io.Reader) (int64, error) {
	file, err := upload.client.OpenFile(upload.binPath, os.O_WRONLY|os.O_APPEND)
	if err != nil {
		return 0, err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(src)
	contents := buf.Bytes()
	n, err := file.Write(contents)
	if err != nil {
		err = fmt.Errorf("Error writing to file: " + err.Error())
	}
	defer file.Close()

	if err == io.ErrUnexpectedEOF {
		err = nil
	}

	m := int64(n)
	upload.info.Offset += m
	return m, err
}

func (upload *sftpUpload) GetReader(ctx context.Context) (io.Reader, error) {
	return upload.client.Open(upload.binPath)
}

func (upload *sftpUpload) Terminate(ctx context.Context) error {
	if err := upload.client.Remove(upload.infoPath); err != nil {
		return err
	}
	if err := upload.client.Remove(upload.binPath); err != nil {
		return err
	}
	return nil
}

func (upload *sftpUpload) ConcatUploads(ctx context.Context, uploads []handler.Upload) (err error) {
	file, err := upload.client.OpenFile(upload.binPath, os.O_WRONLY|os.O_APPEND)
	if err != nil {
		return err
	}
	defer file.Close()

	for range uploads {
		src, err := upload.client.Open(upload.binPath)
		if err != nil {
			return err
		}

		if _, err := io.Copy(file, src); err != nil {
			return err
		}
	}

	return
}

func (upload *sftpUpload) DeclareLength(ctx context.Context, length int64) error {
	upload.info.Size = length
	upload.info.SizeIsDeferred = false
	return upload.writeInfo()
}

// writeInfo updates the entire information. Everything will be overwritten.
func (upload *sftpUpload) writeInfo() error {
	data, err := json.Marshal(upload.info)
	if err != nil {
		return err
	}

	file, err := upload.client.OpenFile(upload.infoPath, os.O_CREATE|os.O_WRONLY)
	file.Write(data)
	defer file.Close()

	return err
}

func (upload *sftpUpload) FinishUpload(ctx context.Context) error {
	//upload.client.Close()
	return nil
}
