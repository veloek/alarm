package daemon

import (
	"bytes"
	"io"
	"io/ioutil"

	mp3 "github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
)

func playSound() error {
	f, err := Asset("alarm.mp3")
	if err != nil {
		return err
	}

	d, err := mp3.NewDecoder(ioutil.NopCloser(bytes.NewReader(f)))
	if err != nil {
		return err
	}
	defer d.Close()

	p, err := oto.NewPlayer(d.SampleRate(), 2, 2, 8192)
	if err != nil {
		return err
	}
	defer p.Close()

	if _, err := io.Copy(p, d); err != nil {
		return err
	}
	return nil
}
