package filetransfer

import (
	"log"
	"os"
	"path/filepath"

	"gioui.org/x/explorer"
	"github.com/claudiodangelis/qrcp/config"
	"github.com/claudiodangelis/qrcp/payload"
	"github.com/claudiodangelis/qrcp/server"
)

func SendFiles(expl *explorer.Explorer) string {
	filepaths := chooseFiles(expl)
	if len(filepaths) == 0 {
		return ""
	}
	downloadLink := createServer(filepaths)
	return downloadLink
}

func ReceiveFiles(receiveUrlChannel, outputDirChannel chan string) {
	cfg, err := config.New("", config.Options{})
	if err != nil {
		log.Fatal(err)
	}
	cfg.Output = getDownloadPath()
	srv, err := server.New(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	if err := srv.ReceiveTo(cfg.Output); err != nil {
		log.Fatal(err)
	}
	receiveUrlChannel <- srv.ReceiveURL
	if err := srv.Wait(); err != nil {
		log.Fatal(err)
	}
	outputDirChannel <- cfg.Output
}

func chooseFiles(expl *explorer.Explorer) []string {
	files, err := expl.ChooseFiles()
	if err != nil {
		if err == explorer.ErrUserDecline {
			return nil
		}
		log.Fatal("failed opening file: ", err)
	}
	var fileInterface interface{}
	filepaths := make([]string, 0)
	for _, currentFile := range files {
		fileInterface = currentFile
		file, ok := fileInterface.(*os.File)
		if !ok {
			log.Fatal("failed opening files")
		}
		filepaths = append(filepaths, file.Name())
	}
	return filepaths
}

func createServer(filepaths []string) string {
	payload, err := payload.FromArgs(filepaths, false)
	if err != nil {
		log.Fatal(err)
	}
	cfg, err := config.New("", config.Options{})
	if err != nil {
		log.Fatal(err)
	}
	srv, err := server.New(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	srv.Send(payload)
	return srv.SendURL
}

func getDownloadPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	downloadDir := filepath.Join(homeDir, "Downloads")
	if _, err := os.Stat(downloadDir); err == nil {
		return downloadDir
	}
	if _, err := os.Stat(homeDir); err == nil {
		return homeDir
	}
	return ""
}
