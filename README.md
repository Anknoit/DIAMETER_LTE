# DIAMETER_LTE
DIAMETER implementation for LTE Applications using RFC 6733

1. Use cases of DIAMETER
   - Wi-Fi authentication (NAS → AAA)

   - Mobile data charging (PCEF → OCS)

   - Broadband PPPoE auth

   - API access control

   - IoT device auth

## ERRORS
2025-12-15.11:24:02.845|A|START log
2025-12-15.11:24:02.845|E|unknown name [CCR] for [command]
2025-12-15.11:24:02.845|E|Error while creating message
2025-12-15.11:24:02.845|E|Bad scenario definition
2025-12-15.11:24:02.845|E|Traffic scenario error
2025-12-15.11:24:02.846|A|STOP  log

start_cleitn.ksh
export LD_LIBRARY_PATH=/usr/local/bin
seagull -conf ../config/conf.client.xml -dico ../config/base_cx.xml -scen ../scenario/ccr-cca.client.xml -log ../logs/ccr-cca.client.log -llevel ET

Fix - 
since CCR-CCA is not base DIAMETER protocol commands and belong to RFC 4006, load -dico config base_cc.xml (cc = credit control)
