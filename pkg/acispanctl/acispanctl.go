// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package acispanctl

import (
	"errors"
	"fmt"
	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type Epg struct {
	Name   string `mapstructure:"name"`
	Tenant string `mapstructure:"tenant"`
	Ap     string `mapstructure:"ap"`
}

type Path struct {
	Pod  string `mapstructure:"pod"`
	Node string `mapstructure:"node"`
	Port string `mapstructure:"port"`
}

type CEP struct {
	Tenant string `mapstructure:"tenant"`
	Ap     string `mapstructure:"ap"`
	Epg    string `mapstructure:"epg"`
	Mac    string `mapstructure:"mac"`
}

type Destination struct {
	Name   string `mapstructure:"name"`
	Ip     string `mapstructure:"ip"`
	Flowid uint   `mapstructure:"flowid"`
	Ttl    uint8  `mapstructure:"ttl"`
	Mtu    uint   `mapstructure:"mtu"`
	Dscp   string `mapstructure:"dscp"`
}

type DestinationGroup struct {
	Name         string        `mapstructure:"name"`
	Tag          string        `mapstructure:"tag"`
	Destinations []Destination `mapstructure:"destinations"`
}

type Source struct {
	Name      string `mapstructure:"name"`
	Direction string `mapstructure:"direction"`
	Epg       Epg    `mapstructure:"epg"`
	Path      Path   `mapstructure:"path"`
	Cep       CEP    `mapstructure:"cep"`
}

type SpanSession struct {
	Name               string             `mapstructure:"name"`
	Destination_Groups []DestinationGroup `mapstructure:"destination_groups"`
	Sources            []Source           `mapstructure:"sources"`
	Admin_state        string             `mapstructure:"admin_state"`
	State              string             `mapstructure:"state"`
}

type SpanConfig struct {
	Sessions []SpanSession `mapstructure:"sessions"`
}

func NewSpanCEPSession(prefix string, tn string, ap string, epg string, mac string, destip string) SpanSession {
	//span session
	spanSession := SpanSession{}
	if prefix == "" {
		prefix = "vs"
	}
	vsname := fmt.Sprintf("%s-vspan", prefix)
	spanSession.Name = vsname

	//destination group
	dstGrp := DestinationGroup{}
	dstgrpname := fmt.Sprintf("%s-dstgrp1", vsname)
	dstGrp.Name = dstgrpname
	dstGrp.Tag = "Yellow Green"

	//destination
	dst := Destination{}
	dst.Name = fmt.Sprintf("%s-dest1", dstgrpname)
	dst.Ip = destip
	dst.Flowid = 1
	dst.Ttl = 64
	dst.Mtu = 1518
	dst.Dscp = ""
	//assign destination to destination group
	dstGrp.Destinations = append(dstGrp.Destinations, dst)

	//span source
	spanSource := Source{}
	srcname := fmt.Sprintf("%s-source1", vsname)
	spanSource.Name = srcname
	spanSource.Direction = "Both"

	//CEP
	cep := CEP{}
	cep.Tenant = tn
	cep.Ap = ap
	cep.Epg = epg
	cep.Mac = mac
	//assign cep to span source
	spanSource.Cep = cep

	//assign destination group to span session
	spanSession.Destination_Groups = append(spanSession.Destination_Groups, dstGrp)

	//assign span source to span session
	spanSession.Sources = append(spanSession.Sources, spanSource)
	spanSession.Admin_state = "start"
	spanSession.State = "present"

	return spanSession
}

func SaveSpanConfig(c SpanConfig, filename string) error {
	bytes, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, bytes, 0644)
}

func GetAPICClient() *client.Client {
	hosturl := fmt.Sprintf("%s://%s", viper.GetString("aciprotocol"), viper.GetString("acihost"))
	username := viper.GetString("aciauth.username")
	password := viper.GetString("aciauth.password")

	return client.GetClient(hosturl, username, client.Password(password), client.Insecure(true))
}

//func CreateVSPANSessionX(sgname string, sources []string, desc string, dstips []string, lbltag string, flowid string,
//	ttl string, mtu string, dscp string) (*models.SpanVSrcGrp, error) {
//
//	c := GetAPICClient()
//	var spanDstGrpS []*models.SpanVDestGrp
//
//	//dst group config
//	for i, dstip := range(dstips) {
//		dgname := fmt.Sprintf("%s-dstgrp-%d", sgname,(i+1))
//		spanDstGrp, err := c.ServiceManager.CreateSpanVDestGrp(dgname, desc, models.SpanVDestGrpAttributes{})
//		if err != nil {
//			fmt.Printf("error while creating span destination group %s\n", dgname)
//			os.Exit(1)
//		}
//		//spanDstGrpS[i] = spanDstGrp
//		spanDstGrpS = append(spanDstGrpS, spanDstGrp)
//
//		dstname := fmt.Sprintf("%s-dst-%d", sgname,(i+1))
//		spanDst, err := c.ServiceManager.CreateSpanVDest(dstname, models.GetMOName(spanDstGrp.DistinguishedName), desc, models.SpanVDestAttributes{})
//		if err != nil {
//			fmt.Printf("error while creating span destination %s, spanDstGrp=%s\n",dstname, models.GetMOName(spanDstGrp.DistinguishedName))
//			os.Exit(1)
//		}
//		spanVEpgSummaryattr := models.SpanVEpgSummaryAttributes{}
//		spanVEpgSummaryattr.DstIp = dstip
//		_, err = c.ServiceManager.CreateSpanVEpgSummary(models.GetMOName(spanDstGrp.DistinguishedName),
//			models.GetMOName(spanDst.DistinguishedName), desc, spanVEpgSummaryattr)
//		if err != nil {
//			fmt.Println("error while creating span destination EPG")
//			os.Exit(1)
//		}
//	}
//
//	//src group config
//	spanSrcGrp, err := c.ServiceManager.CreateSpanVSrcGrp(sgname, desc, models.SpanVSrcGrpAttributes{})
//	if err != nil {
//		fmt.Println("error while creating span source group")
//		os.Exit(1)
//	}
//
//	for _, spanDstGrp := range(spanDstGrpS) {
//		dgname := models.GetMOName(spanDstGrp.DistinguishedName)
//		spanLbl, err := c.ServiceManager.CreateSpanVSpanLbl(dgname, models.GetMOName(spanSrcGrp.DistinguishedName), desc, models.SpanVSpanLblAttributes{})
//		if err != nil {
//			fmt.Printf("error while creating span label %s, spanSrcGrp=%s\n",dgname, models.GetMOName(spanSrcGrp.DistinguishedName))
//			os.Exit(1)
//		}
//		_ = spanLbl
//	}
//
//	for i, source := range(sources) {
//		srcname := fmt.Sprintf("%s-source-%d-%s", sgname,(i+1), source)
//		spanVSrc, err := c.ServiceManager.CreateSpanVSrc(srcname, models.GetMOName(spanSrcGrp.DistinguishedName), desc, models.SpanVSrcAttributes{})
//		if err != nil {
//			fmt.Printf("error while creating span source %s in spanSrcGrp=%s\n",srcname, models.GetMOName(spanSrcGrp.DistinguishedName))
//			os.Exit(1)
//		}
//		_ = spanVSrc
//	}
//	return spanSrcGrp, err
//}

func StartVSPANSession(spanSrcGrpName string) error {
	c := GetAPICClient()
	_, err := c.ServiceManager.StateUpdateSpanVSrcGrp(spanSrcGrpName, "start")
	return err
}

func StopVSPANSession(spanSrcGrpName string) error {
	c := GetAPICClient()
	spanSrcGrp, err := c.ServiceManager.ReadSpanVSrcGrp(spanSrcGrpName)
	if err != nil {
		fmt.Printf("error while reading span group %s\n", spanSrcGrpName)
		os.Exit(1)
	}
	spanVSrcGrpAttr := spanSrcGrp.SpanVSrcGrpAttributes
	desc := spanSrcGrp.Description
	spanVSrcGrpAttr.AdminSt = "stop"
	_, err = c.UpdateSpanVSrcGrp(spanSrcGrpName, desc, spanVSrcGrpAttr)
	if err != nil {
		fmt.Printf("error while updating span group %s\n", spanSrcGrpName)
		os.Exit(1)
	}
	return nil
}

func ApplyVSPANConfig(config SpanConfig) error {
	c := GetAPICClient()
	//fmt.Println("creating span source group")
	desc := "created using Go binding lib" // FIXME: read from config
	for _, session := range config.Sessions {
		// Create span session group
		spanVSrcGrpAttr := models.SpanVSrcGrpAttributes{}
		spanVSrcGrpAttr.AdminSt = session.Admin_state
		//fmt.Printf("Session: %s, state: %s\n", session.Name, session.Admin_state)
		sessionState := strings.ToLower(session.State) // present(add) or absent(delete)
		if sessionState == "present" {
			spanSrcGrp, err := c.ServiceManager.CreateSpanVSrcGrp(session.Name, desc, spanVSrcGrpAttr)
			if err != nil {
				fmt.Printf("error while creating span source group: %s\n", session.Name)
				return err
			}
			_ = spanSrcGrp

			// Create span destination groups
			for _, destg := range session.Destination_Groups {
				dgname := destg.Name
				spanDstGrp, err := c.ServiceManager.CreateSpanVDestGrp(dgname, desc, models.SpanVDestGrpAttributes{})
				if err != nil {
					fmt.Printf("error while creating span destination group %s\n", dgname)
					return err
				}
				for _, dest := range destg.Destinations {
					dstname := dest.Name
					spanDst, err := c.ServiceManager.CreateSpanVDest(dstname, models.GetMOName(spanDstGrp.DistinguishedName), desc, models.SpanVDestAttributes{})
					if err != nil {
						fmt.Printf("error while creating span destination %s, spanDstGrp=%s\n", dstname, models.GetMOName(spanDstGrp.DistinguishedName))
						return err
					}
					spanVEpgSummaryattr := models.SpanVEpgSummaryAttributes{}
					spanVEpgSummaryattr.DstIp = dest.Ip
					_, err = c.ServiceManager.CreateSpanVEpgSummary(models.GetMOName(spanDstGrp.DistinguishedName),
						models.GetMOName(spanDst.DistinguishedName), desc, spanVEpgSummaryattr)
					if err != nil {
						fmt.Println("error while creating span destination EPG")
						return err
					}
				}
			}

			// Create source span labels for every destination groups defined in the config
			for _, destg := range session.Destination_Groups {
				dgname := destg.Name
				_, err := c.ServiceManager.CreateSpanVSpanLbl(dgname, models.GetMOName(spanSrcGrp.DistinguishedName), desc, models.SpanVSpanLblAttributes{})
				if err != nil {
					fmt.Printf("error while creating span label %s, spanSrcGrp=%s\n", dgname, models.GetMOName(spanSrcGrp.DistinguishedName))
					return err
				}
			}

			// Create span sources
			for _, source := range session.Sources {
				srcname := source.Name
				spanVSrcattr := models.SpanVSrcAttributes{}
				spanVSrcattr.Dir = strings.ToLower(source.Direction)
				_, err := c.ServiceManager.CreateSpanVSrc(srcname, models.GetMOName(spanSrcGrp.DistinguishedName), desc, spanVSrcattr)
				if err != nil {
					fmt.Printf("error while creating span source %s in spanSrcGrp=%s\n", srcname, models.GetMOName(spanSrcGrp.DistinguishedName))
					return err
				}
				//parentDn is common for all 3 types
				parentDn := fmt.Sprintf("uni/infra/vsrcgrp-%s/vsrc-%s", models.GetMOName(spanSrcGrp.DistinguishedName), source.Name)
				// PathEp
				if source.Path.Port != "" && source.Path.Node != "" && source.Path.Pod != "" {
					tDn := fmt.Sprintf("topology/pod-%s/paths-%s/pathep-[%s]", source.Path.Pod, source.Path.Node, source.Path.Port)
					err = c.ServiceManager.CreateRelationSpanRsSrcToPathEp(parentDn, tDn)
					if err != nil {
						fmt.Printf("error while creating span source path %s in spanSrcGrp=%s\n", srcname, models.GetMOName(spanSrcGrp.DistinguishedName))
						return err
					}
				}
				// VPort
				if source.Cep.Mac != "" && source.Cep.Epg != "" && source.Cep.Ap != "" && source.Cep.Tenant != "" {
					tDn := fmt.Sprintf("uni/tn-%s/ap-%s/epg-%s/cep-%s", source.Cep.Tenant, source.Cep.Ap, source.Cep.Epg, source.Cep.Mac)
					err = c.ServiceManager.CreateRelationSpanRsSrcToVPort(parentDn, tDn)
					if err != nil {
						fmt.Printf("error while creating span source CEP %s in spanSrcGrp=%s\n", srcname, models.GetMOName(spanSrcGrp.DistinguishedName))
						return err
					}
				}
				// Epg
				if source.Epg.Name != "" && source.Epg.Ap != "" && source.Epg.Tenant != "" {
					tDn := fmt.Sprintf("uni/tn-%s/ap-%s/epg-%s", source.Epg.Tenant, source.Epg.Ap, source.Epg.Name)
					err = c.ServiceManager.CreateRelationSpanRsSrcToEpg(parentDn, tDn)
					if err != nil {
						fmt.Printf("error while creating span source epg %s in spanSrcGrp=%s\n", srcname, models.GetMOName(spanSrcGrp.DistinguishedName))
						return err
					}
				}
			}
		}
		if sessionState == "absent" {
			// delete the span session
			//fmt.Printf("deleting span source group: %s", session.Name)
			err := c.ServiceManager.DeleteSpanVSrcGrp(session.Name)
			if err != nil {
				fmt.Printf("error while deleting span source group: %s", session.Name)
				return err
			}

			// delete all related span destination groups
			for _, destg := range session.Destination_Groups {
				dgname := destg.Name
				err := c.ServiceManager.DeleteSpanVDestGrp(dgname)
				if err != nil {
					fmt.Printf("error while deleting span destination group %s\n", dgname)
					return err
				}
			}
		}
	}
	return nil
}

func DeleteVSPANSession(sgname string) error {
	c := GetAPICClient()

	spanLblMol, err := c.ServiceManager.ListSpanVSpanLbl()
	if err != nil {
		//attempt  again
		for attempt := 0; attempt < 5; attempt++ {
			spanLblMol, err = c.ServiceManager.ListSpanVSpanLbl()
			if err == nil || attempt == 5 {
				break
			}
			time.Sleep(time.Second * 1)
		}
	}
	for _, spanLblMo := range spanLblMol {

		if len(strings.Split(spanLblMo.DistinguishedName, "/")) != 4 {
			continue
		}
		srcgrpname := strings.TrimPrefix(strings.Split(spanLblMo.DistinguishedName, "/")[2], "vsrcgrp-")
		if srcgrpname == sgname {
			//fmt.Printf("deleting...%s\n", spanLblMo.Name)
			err := c.ServiceManager.DeleteSpanVDestGrp(spanLblMo.Name)
			maxAttempted := false
			if err != nil {
				//attempt  again
				for attempt := 0; attempt < 5; attempt++ {
					err = c.ServiceManager.DeleteSpanVDestGrp(spanLblMo.Name)
					if err == nil || attempt == 5 {
						if attempt == 5 {
							maxAttempted = true
						}
						break
					}
					time.Sleep(time.Second * 1)
				}
			}

			if err != nil || maxAttempted {
				fmt.Printf("error deleting...%s\n", spanLblMo.Name)
				fmt.Println(err)
			}
		}
	}
	return c.ServiceManager.DeleteSpanVSrcGrp(sgname)
}

func PrintAllOpflexIDEp() {
	c := GetAPICClient()
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	opflexIDEps, err := c.ServiceManager.ListOpflexIDEp()
	if err != nil {
		fmt.Println("error while retrieving all opflex IDEps")
		fmt.Println(err)
		os.Exit(1)
	}
	table.SetHeader([]string{"Container", "Domain Name", "Tenant/AP/EPG", "IP", "MAC"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetRowLine(true)
	table.SetRowSeparator("-")

	for _, opflexIDEp := range opflexIDEps {
		if opflexIDEp.ContainerName != "" {
			s := strings.Split(opflexIDEp.EpgPKey, "/")
			tenant := strings.TrimPrefix(s[1], "tn-")
			ap := strings.TrimPrefix(s[2], "ap-")
			epg := strings.TrimPrefix(s[3], "epg-")
			tae := fmt.Sprintf("%s, %s, %s", tenant, ap, epg)
			//fmt.Printf("Container: %s, EPG: %s, IP: %s, MAC: %s\n", opflexIDEp.ContainerName, opflexIDEp.EpgPKey, opflexIDEp.Ip, opflexIDEp.Mac)
			table.Append([]string{opflexIDEp.ContainerName, opflexIDEp.DomName, tae, opflexIDEp.Ip, opflexIDEp.Mac})
		}
	}
	table.Render()
	fmt.Println(tableString.String())
}

func CreateSpanSessionFromCont(cont string, domain string, namespace string, dst_ip string) error {
	c := GetAPICClient()
	opflexIDEps, err := c.ServiceManager.ListOpflexIDEp()
	if err != nil {
		fmt.Println("error while retrieving all opflex IDEps")
		os.Exit(1)
	}
	spanConfig := SpanConfig{}
	for _, opflexIDEp := range opflexIDEps {
		if opflexIDEp.ContainerName == cont && opflexIDEp.DomName == domain && opflexIDEp.Namespace == namespace {
			s := strings.Split(opflexIDEp.EpgPKey, "/")
			tenant := strings.TrimPrefix(s[1], "tn-")
			ap := strings.TrimPrefix(s[2], "ap-")
			epg := strings.TrimPrefix(s[3], "epg-")
			spanConfig.Sessions = append(spanConfig.Sessions, NewSpanCEPSession(opflexIDEp.ContainerName, tenant, ap, epg, opflexIDEp.Mac, dst_ip))
			fname := fmt.Sprintf("%s-vspan.yaml", cont)
			return SaveSpanConfig(spanConfig, fname)
		}
	}
	//fmt.Printf("container %s not found on ACI\n", cont)
	err_str := fmt.Sprintf("container %s not found on ACI", cont)
	return errors.New(err_str)
}

func PrintAllVSPANSessions() {
	c := GetAPICClient()
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)

	spansessions, err := c.ServiceManager.ListSpanVSrcGrp()
	if err != nil {
		fmt.Println("error while retrieving all erspan sessions")
		os.Exit(1)
	}
	var spanLblMap = make(map[string]string)
	var spanSrcMap = make(map[string]string)
	var spanDstIpMap = make(map[string]string)

	for _, session := range spansessions {
		if models.GetMOName(session.DistinguishedName) == "default" {
			//skip processing default sessions
			continue
		}

		spanlbll := ""
		srctaeml := ""
		spanLblMol, err := c.ServiceManager.ListSpanVSpanLbl()
		if err != nil {
			//attempt  again
			for attempt := 0; attempt < 5; attempt++ {
				spanLblMol, err = c.ServiceManager.ListSpanVSpanLbl()
				if err == nil || attempt == 5 {
					break
				}
				time.Sleep(time.Second * 1)
			}
		}
		//map all span destination group and destinatipon IP for all span sessions
		for _, spanLblMo := range spanLblMol {
			dstipl := ""
			if len(strings.Split(spanLblMo.DistinguishedName, "/")) != 4 {
				continue
			}
			srcgrpname := strings.TrimPrefix(strings.Split(spanLblMo.DistinguishedName, "/")[2], "vsrcgrp-")
			if srcgrpname == models.GetMOName(session.DistinguishedName) {
				spanlbll = fmt.Sprintf("%s %s", spanlbll, models.GetMOName(spanLblMo.DistinguishedName))
				spanDstl, err := c.ServiceManager.ListSpanVDest()
				if err != nil {
					fmt.Printf("error while retrieving erspan destinations for %s\n", models.GetMOName(spanLblMo.DistinguishedName))
					fmt.Println(err)
					os.Exit(1)
				}
				for _, spanDst := range spanDstl {
					if strings.TrimPrefix(strings.Split(spanDst.DistinguishedName, "/")[2], "vdestgrp-") != models.GetMOName(spanLblMo.DistinguishedName) {
						continue
					}
					epg_attr, err := c.ServiceManager.ReadSpanVEpgSummary(models.GetMOName(spanLblMo.DistinguishedName), models.GetMOName(spanDst.DistinguishedName))
					if err != nil {
						isMaxAttempt := false
						for attempt := 0; attempt < 5; attempt++ {
							epg_attr, err = c.ServiceManager.ReadSpanVEpgSummary(models.GetMOName(spanLblMo.DistinguishedName), models.GetMOName(spanDst.DistinguishedName))
							if err == nil || attempt == 5 {
								if attempt == 5 {
									isMaxAttempt = true
								}
								break
							}
							time.Sleep(time.Second * 1)
						}
						if isMaxAttempt {
							fmt.Printf("error while reading SPAN epg summary for %s\n", models.GetMOName(spanLblMo.DistinguishedName))
							os.Exit(1)
						}

					}
					dstipl = fmt.Sprintf("%s %s", dstipl, epg_attr.DstIp)
				}
				spanDstIpMap[models.GetMOName(session.DistinguishedName)] = fmt.Sprintf("%s %s", spanDstIpMap[models.GetMOName(session.DistinguishedName)], dstipl)
			}
		}
		spanLblMap[models.GetMOName(session.DistinguishedName)] = spanlbll

		spanSrcMol, err := c.ServiceManager.ListSpanVSrc(models.GetMOName(session.DistinguishedName))
		if err != nil {
			//attempt again
			for attempt := 0; attempt < 5; attempt++ {
				spanSrcMol, err = c.ServiceManager.ListSpanVSrc(models.GetMOName(session.DistinguishedName))
				if err == nil || attempt == 5 {
					break
				}
				time.Sleep(time.Second * 1)
			}
		}
		for _, spanSrcMo := range spanSrcMol {
			spanRsSrcToVPortL, err := c.ServiceManager.ListSpanRsSrcToVPort()
			if err != nil {
				//attempt again
				isMaxAttempt := false
				for attempt := 0; attempt < 5; attempt++ {
					spanRsSrcToVPortL, err = c.ServiceManager.ListSpanRsSrcToVPort()
					if err == nil || attempt == 5 {
						if attempt == 5 {
							isMaxAttempt = true
						}
						break
					}
					time.Sleep(time.Second * 1)
				}
				if isMaxAttempt {
					fmt.Printf("error while retrieving rs src to vport for %s in %s\n", models.GetMOName(spanSrcMo.DistinguishedName),
						models.GetMOName(session.DistinguishedName))
					os.Exit(1)
				}
			}

			for _, spanRsSrcToVPort := range spanRsSrcToVPortL {
				session_name_str := strings.TrimPrefix(strings.Split(spanRsSrcToVPort.DistinguishedName, "/")[2], "vsrcgrp-")
				span_src_name_str := strings.TrimPrefix(strings.Split(spanRsSrcToVPort.DistinguishedName, "/")[3], "vsrc-")
				if session_name_str != models.GetMOName(session.DistinguishedName) && span_src_name_str != models.GetMOName(spanSrcMo.DistinguishedName) {
					continue
				}
				tdn := spanRsSrcToVPort.Tdn
				s := strings.Split(tdn, "/")
				tenant := strings.TrimPrefix(s[1], "tn-")
				ap := strings.TrimPrefix(s[2], "ap-")
				epg := strings.TrimPrefix(s[3], "epg-")
				mac := strings.TrimPrefix(s[4], "cep-")
				taem := fmt.Sprintf("%s/%s/%s/%s", tenant, ap, epg, mac)
				srctaeml = fmt.Sprintf("%s %s", srctaeml, strings.TrimSpace(taem))

			}

		}
		spanSrcMap[models.GetMOName(session.DistinguishedName)] = strings.TrimSpace(srctaeml)

	}

	table.SetHeader([]string{"SPAN Name", "Tenant/AP/EPG/MAC", "SPAN Destination IP", "Admin State"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetRowLine(true)
	table.SetRowSeparator("-")
	for _, session := range spansessions {
		if models.GetMOName(session.DistinguishedName) != "default" {
			table.Append([]string{models.GetMOName(session.DistinguishedName), spanSrcMap[models.GetMOName(session.DistinguishedName)],
				spanDstIpMap[models.GetMOName(session.DistinguishedName)], session.AdminSt})
		}
	}

	table.Render()
	fmt.Println(tableString.String())
}
