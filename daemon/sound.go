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

	c, err := oto.NewContext(d.SampleRate(), 2, 2, 8192)
	if err != nil {
		return err
	}
	defer c.Close()

	p := c.NewPlayer()
	defer p.Close()

	if _, err := io.Copy(p, d); err != nil {
		return err
	}
	return nil
}
