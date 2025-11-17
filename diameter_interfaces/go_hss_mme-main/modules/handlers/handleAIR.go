package handlers

import (
	"fmt"
	"log"

	"epiphanyHSS/modules/db/mysql"
	"epiphanyHSS/modules/go-diameter/diam"
	"epiphanyHSS/modules/go-diameter/diam/avp"
	"epiphanyHSS/modules/go-diameter/diam/datatype"
	"epiphanyHSS/modules/go-diameter/diam/sm"
	"io"
)

const VENDOR_3GPP = 10415

func HandleAIR(settings sm.Settings) diam.HandlerFunc {
	type RequestedEUTRANAuthInfo struct {
		NumVectors        datatype.Unsigned32  `avp:"Number-Of-Requested-Vectors"`
		ImmediateResponse datatype.Unsigned32  `avp:"Immediate-Response-Preferred"`
		ResyncInfo        datatype.OctetString `avp:"Re-synchronization-Info"`
	}

	type AIR struct {
		SessionID               datatype.UTF8String       `avp:"Session-Id"`
		OriginHost              datatype.DiameterIdentity `avp:"Origin-Host"`
		OriginRealm             datatype.DiameterIdentity `avp:"Origin-Realm"`
		AuthSessionState        datatype.UTF8String       `avp:"Auth-Session-State"`
		UserName                string                    `avp:"User-Name"`
		VisitedPLMNID           datatype.OctetString      `avp:"Visited-PLMN-Id"`
		RequestedEUTRANAuthInfo RequestedEUTRANAuthInfo   `avp:"Requested-EUTRAN-Authentication-Info"`
	}
	return func(c diam.Conn, m *diam.Message) {
		var err error
		var req AIR
		var code uint32

		err = m.Unmarshal(&req)
		if err != nil {
			err = fmt.Errorf("Unmarshal failed: %s", err)
			code = diam.UnableToComply
			log.Printf("Invalid AIR(%d): %s\n", code, err.Error())
		} else {
			userProfile, err := mysql.FetchRDSUserProfile(req.UserName)
			if err != nil {
				fmt.Println("Error : %v", err)
				code = diam.UnableToComply
			}
			if userProfile == nil {
				fmt.Println("user not found")
				code = diam.UnableToComply
			} else {
				fmt.Println("user found")
				code = diam.Success
			}
		}

		a := m.Answer(code)
		a.InsertAVP(diam.NewAVP(avp.SessionID, avp.Mbit, 0, req.SessionID))
		a.NewAVP(avp.OriginHost, avp.Mbit, 0, settings.OriginHost)
		a.NewAVP(avp.OriginRealm, avp.Mbit, 0, settings.OriginRealm)
		a.NewAVP(avp.OriginStateID, avp.Mbit, 0, settings.OriginStateID)
		_, err = sendAIA(settings, c, a)
		if err != nil {
			log.Printf("Failed to send AIA: %s", err.Error())
		}
	}
}

func sendAIA(settings sm.Settings, w io.Writer, m *diam.Message) (n int64, err error) {
	m.NewAVP(avp.AuthenticationInfo, avp.Mbit, VENDOR_3GPP, &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.EUTRANVector, avp.Mbit, VENDOR_3GPP, &diam.GroupedAVP{
				AVP: []*diam.AVP{
					diam.NewAVP(avp.RAND, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.OctetString("\x94\xbf/T\xc3v\xf3\x0e\x87\x83\x06k'\x18Z\x19")),
					diam.NewAVP(avp.XRES, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.OctetString("F\xf0\"\xb9%#\xf58")),
					diam.NewAVP(avp.AUTN, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.OctetString("\xc7G!;\xad~\x80\x00)\x08o%\x11\x0cP_")),
					diam.NewAVP(avp.KASME, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.OctetString("\xbf\x00\xf9\x80h3\"\x0e\xa1\x1c\xfa\x93\x03@\xd6\xf8\x02\xd51Y\xeb\xc4\x9d=\t\x14{\xeb!\xec\xcb:")),
				},
			}),
		},
	})

	return m.WriteTo(w)
}
