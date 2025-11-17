package handlers

import (
	"fmt"
	"log"

	"epiphanyHSS/modules/go-diameter/diam"
	"epiphanyHSS/modules/go-diameter/diam/avp"
	"epiphanyHSS/modules/go-diameter/diam/datatype"
	"epiphanyHSS/modules/go-diameter/diam/sm"
	"io"
)

func HandleULR(settings sm.Settings) diam.HandlerFunc {
	type ULR struct {
		SessionID        datatype.UTF8String       `avp:"Session-Id"`
		OriginHost       datatype.DiameterIdentity `avp:"Origin-Host"`
		OriginRealm      datatype.DiameterIdentity `avp:"Origin-Realm"`
		AuthSessionState datatype.Unsigned32       `avp:"Auth-Session-State"`
		UserName         datatype.UTF8String       `avp:"User-Name"`
		VisitedPLMNID    datatype.OctetString      `avp:"Visited-PLMN-Id"`
		RATType          datatype.Unsigned32       `avp:"RAT-Type"`
		ULRFlags         datatype.Unsigned32       `avp:"ULR-Flags"`
	}
	return func(c diam.Conn, m *diam.Message) {
		var err error = nil
		var req ULR
		var code uint32

		err = m.Unmarshal(&req)
		if err != nil {
			err = fmt.Errorf("Unmarshal failed: %s", err)
			code = diam.UnableToComply
			log.Printf("Invalid ULR(%d): %s\n", code, err.Error())
		} else {
			code = diam.Success
		}

		a := m.Answer(code)
		a.InsertAVP(diam.NewAVP(avp.SessionID, avp.Mbit, 0, req.SessionID))
		a.NewAVP(avp.AuthSessionState, avp.Mbit, 0, req.AuthSessionState)
		a.NewAVP(avp.OriginHost, avp.Mbit, 0, settings.OriginHost)
		a.NewAVP(avp.OriginRealm, avp.Mbit, 0, settings.OriginRealm)
		a.NewAVP(avp.OriginStateID, avp.Mbit, 0, settings.OriginStateID)
		_, err = sendULA(settings, c, a)
		if err != nil {
			log.Printf("Failed to send ULA: %s", err.Error())
		}
	}
}

func sendULA(settings sm.Settings, w io.Writer, m *diam.Message) (n int64, err error) {
	m.NewAVP(avp.ULAFlags, avp.Mbit|avp.Vbit, VENDOR_3GPP, datatype.Unsigned32(1))
	// Add other AVPs as necessary
	return m.WriteTo(w)
}
