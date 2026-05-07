package recorders

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type FileRecorder struct {
	baseFilename string
	iqFile       *os.File
	audioFile    *os.File
	dataFile     *os.File
	audioAsWav   bool

	audioFilename string
	iqFilename    string
	dataFilename  string
}

func (f *FileRecorder) Open(config RecorderConfig) bool {
	log.Println("FileRecorder: Writing Metadata to", fmt.Sprintf("%s-metadata.json", config.BaseFilename))
	metaFile, err := os.Create(fmt.Sprintf("%s-metadata.json", config.BaseFilename))

	if err != nil {
		panic(err)
	}

	metadataJson, err := json.MarshalIndent(config.Metadata, "", "   ")

	if err != nil {
		panic(err)
	}

	_, _ = metaFile.Write(metadataJson)

	if err != nil {
		panic(err)
	}

	_ = metaFile.Close()

	f.baseFilename = config.BaseFilename
	f.audioAsWav = config.AudioAsWav
	f.iqFilename = fmt.Sprintf("%s-iq.cfile", f.baseFilename)
	f.audioFilename = fmt.Sprintf("%s-audio.float32", f.baseFilename)
	f.dataFilename = fmt.Sprintf("%s-data.bytes", f.baseFilename)

	return true
}

func (f *FileRecorder) Close() bool {
	if f.audioFile != nil {
		log.Println("FileRecorder: Closing Audio File", f.audioFilename)
		_ = f.audioFile.Close()
		f.audioFile = nil
	}
	if f.iqFile != nil {
		log.Println("FileRecorder: Closing IQ File", f.iqFile)
		_ = f.iqFile.Close()
		f.iqFile = nil
	}
	if f.dataFile != nil {
		log.Println("FileRecorder: Closing Data File", f.dataFile)
		_ = f.dataFile.Close()
		f.dataFile = nil
	}

	return true
}

func (f *FileRecorder) WriteIQ(data []complex64) {
	if f.iqFile == nil {
		var err error
		log.Println("FileRecorder: Writing IQ to", f.iqFilename)
		f.iqFile, err = os.Create(f.iqFilename)
		if err != nil {
			panic(err)
		}
	}

	err := binary.Write(f.iqFile, binary.LittleEndian, data)
	if err != nil {
		panic(err)
	}
}

func (f *FileRecorder) WriteAudio(data []float32) {
	if f.audioFile == nil {
		var err error
		// TODO: Audio as Wave
		log.Println("FileRecorder: Writing Audio to", f.audioFilename)
		f.audioFile, err = os.Create(f.audioFilename)
		if err != nil {
			panic(err)
		}
	}

	err := binary.Write(f.audioFile, binary.LittleEndian, data)
	if err != nil {
		panic(err)
	}
}

func (f *FileRecorder) WriteData(data []byte) {
	if f.dataFile == nil {
		var err error
		log.Println("FileRecorder: Writing Bytes to", f.dataFilename)
		f.dataFile, err = os.Create(f.dataFilename)
		if err != nil {
			panic(err)
		}
	}

	_, err := f.dataFile.Write(data)
	if err != nil {
		panic(err)
	}
}
