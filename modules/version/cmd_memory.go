package version

import (
	"github.com/bwmarrin/discordgo"
	"github.com/innovandalism/shodan"
	"github.com/innovandalism/shodan/util"
	"runtime"
	"strconv"
)

type MemoryCommand struct{}

func (command *MemoryCommand) GetNames() []string {
	return []string{"memory", "mem"}
}

func (command *MemoryCommand) Invoke(ci *shodan.CommandInvocation) bool {
	runtime.ReadMemStats(mod.memStats)
	me := &discordgo.MessageEmbed{
		Type:        "rich",
		Title:       "SHODAN",
		Description: "Memory Statistics",
		Fields: []*discordgo.MessageEmbedField{
			// General Statistics
			{"Alloc", strconv.Itoa(int(mod.memStats.Alloc)), true},
			{"TotalAlloc", strconv.Itoa(int(mod.memStats.TotalAlloc)), true},
			{"Sys", strconv.Itoa(int(mod.memStats.Sys)), true},
			{"Lookups", strconv.Itoa(int(mod.memStats.Lookups)), true},
			{"Mallocs", strconv.Itoa(int(mod.memStats.Mallocs)), true},
			{"Frees", strconv.Itoa(int(mod.memStats.Frees)), true},

			// Heap
			{"HeapAlloc", strconv.Itoa(int(mod.memStats.HeapAlloc)), true},
			{"HeapSys", strconv.Itoa(int(mod.memStats.HeapSys)), true},
			{"HeapIdle", strconv.Itoa(int(mod.memStats.HeapIdle)), true},
			{"HeapInuse", strconv.Itoa(int(mod.memStats.HeapInuse)), true},
			{"HeapObjects", strconv.Itoa(int(mod.memStats.HeapObjects)), true},
			{"HeapReleased", strconv.Itoa(int(mod.memStats.HeapReleased)), true},

			// Stack
			{"StackInuse", strconv.Itoa(int(mod.memStats.StackInuse)), true},
			{"StackSys", strconv.Itoa(int(mod.memStats.StackSys)), true},
			{"MSpanInuse", strconv.Itoa(int(mod.memStats.MSpanInuse)), true},
			{"MSpanSys", strconv.Itoa(int(mod.memStats.MSpanSys)), true},
			{"MCacheInuse", strconv.Itoa(int(mod.memStats.MCacheInuse)), true},
			{"MCacheSys", strconv.Itoa(int(mod.memStats.MCacheSys)), true},
			{"BuckHashSys", strconv.Itoa(int(mod.memStats.BuckHashSys)), true},
			{"GCSys", strconv.Itoa(int(mod.memStats.GCSys)), true},
			{"OtherSys", strconv.Itoa(int(mod.memStats.OtherSys)), true},

			//GC
			{"NumGC", strconv.Itoa(int(mod.memStats.NumGC)), true},
			{"GCCPUFraction", strconv.Itoa(int(mod.memStats.GCCPUFraction)), true},
			{"EnableGC", strconv.FormatBool(mod.memStats.EnableGC), true},
			{"DebugGC", strconv.FormatBool(mod.memStats.DebugGC), true},
		},
	}
	_, err := ci.Session.ChannelMessageSendEmbed(ci.Event.ChannelID, me)
	if err != nil {
		util.ReportThreadError(false, err)
	}
	return true
}
