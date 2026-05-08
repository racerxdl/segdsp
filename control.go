package main

import (
	"encoding/json"
	"fmt"
	"github.com/racerxdl/segdsp/eventmanager"
	"log"
)

func handleControl(msg controlMessage) {
	if !webCanControl {
		return
	}

	changed := false
	rebuild := false

	if msg.ChannelFrequency != nil {
		channelFrequency = *msg.ChannelFrequency
		if radioClient != nil {
			radioClient.SetCenterFrequency(uint32(channelFrequency))
		}
		changed = true
	}

	if msg.FFTfrequency != nil {
		displayFrequency = *msg.FFTfrequency
		if radioClient != nil {
			radioClient.SetSmartCenterFrequency(uint32(displayFrequency))
		}
		changed = true
	}

	if msg.DemodulatorMode != nil {
		for _, m := range modes {
			if m == *msg.DemodulatorMode {
				demodulatorMode = *msg.DemodulatorMode
				rebuild = true
				break
			}
		}
	}

	if msg.FilterBandwidth != nil {
		filterBandwidth = *msg.FilterBandwidth
		rebuild = true
	}

	if msg.Squelch != nil {
		squelch = *msg.Squelch
		rebuild = true
	}

	if msg.DisplayOffset != nil {
		displayOffset = int(*msg.DisplayOffset)
		changed = true
	}

	if msg.DisplayRange != nil {
		displayRange = int(*msg.DisplayRange)
		changed = true
	}

	if msg.Rebuild != nil && *msg.Rebuild {
		rebuild = true
	}

	if rebuild && radioClient != nil {
		sampleRate := radioClient.GetSampleRate()
		dspPipeline.Stop()

		ev := &eventmanager.EventManager{}
		ev.AddHandler(eventmanager.EvSquelchOn, squelchOn)
		ev.AddHandler(eventmanager.EvSquelchOff, squelchOff)

		dspPipeline = NewDSPPipeline()
		dspPipeline.Demodulator = buildDSP(sampleRate)
		dspPipeline.Demodulator.SetEventManager(ev)
		dspPipeline.SetCallback(sendData)
		dspPipeline.Start()

		changed = true
		log.Printf("Pipeline rebuilt: mode=%s bw=%d squelch=%.1f", demodulatorMode, filterBandwidth, squelch)
	}

	if changed {
		onDeviceSync(radioClient)

		resp, _ := json.Marshal(map[string]string{"MessageType": "controlAck"})
		go wsServer.BroadcastMessage(string(resp))
	} else {
		resp, _ := json.Marshal(map[string]string{"MessageType": "controlAck"})
		go wsServer.BroadcastMessage(string(resp))
	}

	_ = fmt.Sprintf("control: %+v", msg)
}
