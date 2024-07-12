package dictionary

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	data := `
    <configuration>
        <subtype id='Enum5GDeliveryReportRequested'>
            <datatype>unsigned int16</datatype>
            <value id='1'>YES</value>
            <value id='2'>NO</value>
        </subtype>
        <container id='MtxTxnActionData'>
            <doc_description>Abstract base container for all transaction actions.  Only actions which derive from this container should appear in the TxnMsg's TxnActionList.</doc_description>
            <key>142</key>
            <created_schema_version>4300</created_schema_version>
            <field id='ObjectId'>
                <doc_description></doc_description>
                <datatype>object id</datatype>
                <created_schema_version>4300</created_schema_version>
            </field>
            <field id='Operator'>
                <doc_description></doc_description>
                <datatype>unsigned int32</datatype>
                <created_schema_version>4300</created_schema_version>
            </field>
            <field id='ChangeCounterAfter'>
                <doc_description>Audit trail: object's change counter after the action has been applied.  Set by transaction server.</doc_description>
                <datatype>unsigned int32</datatype>
                <created_schema_version>4300</created_schema_version>
            </field>
        </container>
        <container id='MtxResponseCreate'>
            <doc_description>Response to creating a new DB object.</doc_description>
            <key>243</key>
            <created_schema_version>4300</created_schema_version>
            <base_container id='4300'>MtxResponse</base_container>
            <field id='ObjectId'>
                <doc_description>The ObjectId for the created object.</doc_description>
                <datatype>object id</datatype>
                <created_schema_version>4300</created_schema_version>
            </field>
        </container>
        <container id='MtxDiamMsg'>
            <doc_description>Common Base Container for all Diameter messages.</doc_description>
            <key>94</key>
            <created_schema_version>4300</created_schema_version>
            <base_container id='4300'>MtxChrgMsg</base_container>
            <field id='ProxyFlag'>
                <doc_description>Diameter Proxiable bit.</doc_description>
                <datatype>bool</datatype>
                <created_schema_version>4300</created_schema_version>
            </field>
            <field id='ApplicationId'>
                <doc_description>Diameter Application-ID AVP.  This is four octets and is used to identify to which application the message is applicable for.  The application can be an authentication application, an accounting application or a vendor specific application.  See RFC-3588 Section 11.3 for the possible values that the application-id may use.  The application-id in the header MUST be the same as what is contained in any relevant AVPs contained in the message.  From RFC-3588 Section 11.3: There are standards-track application ids and vendor specific application ids.  IANA [IANA] has assigned the range 0x00000001 to 0x00ffffff for standards-track applications; and 0x01000000 - 0xfffffffe for vendor specific applications, on a first-come, first-served basis.  The following values are allocated.  Diameter Common Messages 0 NASREQ 1 [NASREQ] Mobile-IP 2 [DIAMMIP] Diameter Base Accounting 3 Credit-Control 4 Relay 0xffffffff Assignment of standards-track application Ids are by Designated Expert with Specification Required [IANA].  Both Application-Id and Acct-Application-Id AVPs use the same Application Identifier space.  Vendor-Specific Application Identifiers, are for Private Use.  Vendor-Specific Application Identifiers are assigned on a First Come, First Served basis by IANA.</doc_description>
                <datatype>unsigned int32</datatype>
                <created_schema_version>4300</created_schema_version>
            </field>
            <field id='HopByHopId'>
                <doc_description>Diameter Hop-by-Hop Identifier AVP.  This is an unsigned 32-bit integer field (in network byte order) and aids in matching requests and replies.  The sender MUST ensure that the Hop-by-Hop identifier in a request is unique on a given connection at any given time, and MAY attempt to ensure that the number is unique across reboots.  The sender of an Answer message MUST ensure that the Hop-by-Hop Identifier field contains the same value that was found in the corresponding request.  The Hop-by-Hop identifier is normally a monotonically increasing number, whose start value was randomly generated.  An answer message that is received with an unknown Hop-by-Hop Identifier MUST be discarded.</doc_description>
                <datatype>unsigned int32</datatype>
                <created_schema_version>4300</created_schema_version>
                <deleted_schema_version>4310</deleted_schema_version>
            </field>
            <field id='EndToEndId'>
                <doc_description>Diameter End-to-End Identifier AVP.  This is an unsigned 32-bit integer field (in network byte order) and is used to detect duplicate messages.  Upon reboot implementations MAY set the high order 12 bits to contain the low order 12 bits of current time, and the low order 20 bits to a random value.  Senders of request messages MUST insert a unique identifier on each message.  The identifier MUST remain locally unique for a period of at least 4 minutes, even across reboots.  The originator of an Answer message MUST ensure that the End-to-End Identifier field contains the same value that was found in the corresponding request.  The End-to-End Identifier MUST NOT be modified by Diameter agents of any kind.  The combination of the Origin-Host (see RFC-3588 Section 6.3) and this field is used to detect duplicates.  Duplicate requests SHOULD cause the same answer to be transmitted (modulo the hop-by-hop Identifier field and any routing AVPs that may be present), and MUST NOT affect any state that was set when the original request was processed.  Duplicate answer messages that are to be locally consumed (see RFC-3588 Section 6.2) SHOULD be silently discarded.</doc_description>
                <datatype>unsigned int32</datatype>
                <created_schema_version>4300</created_schema_version>
                <deleted_schema_version>4310</deleted_schema_version>
            </field>
            <field id='DiamOp'>
                <doc_description>Diameter Command-Code AVP.  Every Diameter message MUST contain a command code in its header's Command-Code field, which is used to determine the action that is to be taken for a particular message.</doc_description>
                <datatype>unsigned int32</datatype>
                <created_schema_version>4300</created_schema_version>
            </field>
            <field id='SessionId'>
                <doc_description>Diameter Session-Id AVP.  This is used to identify a specific session (see RFC-3588 Section 8).  All messages pertaining to a specific session MUST include only one Session-Id AVP and the same value MUST be used throughout the life of a session.  When present, the Session-Id SHOULD appear immediately following the Diameter Header (see RFC-3588 Section 3).  The Session-Id MUST be globally and eternally unique, as it is meant to uniquely identify a user session without reference to any other information, and may be needed to correlate historical authentication information with accounting information.  The Session-Id includes a mandatory portion and an implementation-defined portion; a recommended format for the implementation-defined portion is outlined below.  The Session-Id MUST begin with the sender's identity encoded in the DiameterIdentity type (see RFC-3588 Section 4.4).  The remainder of the Session-Id is delimited by a ";" character, and MAY be any sequence that the client can guarantee to be eternally unique; however, the following format is recommended, (square brackets [] indicate an optional element): &lt;DiameterIdentity&gt;;&lt;high 32 bits&gt;;&lt;low 32 bits&gt;[;&lt;optional value&gt;] &lt;high 32 bits&gt; and &lt;low 32 bits&gt; are decimal representations of the high and low 32 bits of a monotonically increasing 64-bit value.  The 64-bit value is rendered in two part to simplify formatting by 32-bit processors.  At startup, the high 32 bits of the 64-bit value MAY be initialized to the time, and the low 32 bits MAY be initialized to zero.  This will for practical purposes eliminate the possibility of overlapping Session-Ids after a reboot, assuming the reboot process takes longer than a second.  Alternatively, an implementation MAY keep track of the increasing value in non-volatile memory.  &lt;optional value&gt; is implementation specific but may include a modem's device Id, a layer 2 address, timestamp, etc.  Example, in which there is no optional value: accesspoint7.acme.com;1876543210;523 Example, in which there is an optional value: accesspoint7.acme.com;1876543210;523;mobile@200.1.1.88 The Session-Id is created by the Diameter application initiating the session, which in most cases is done by the client.  Note that a Session-Id MAY be used for both the authorization and accounting commands of a given application.</doc_description>
                <datatype>string</datatype>
                <created_schema_version>4300</created_schema_version>
            </field>
            <field id='MultiSessionId'>
                <doc_description>Diameter Acct-Multi-Session-Id AVP.  This is used to link together multiple related accounting sessions, where each session would have a unique Session-Id, but the same Acct-Multi-Session-Id AVP.  This AVP MAY be returned by the Diameter server in an authorization answer, and MUST be used in all accounting messages for the given session.</doc_description>
                <datatype>string</datatype>
                <created_schema_version>4500</created_schema_version>
            </field>
            <field id='SubSessionId'>
                <doc_description>Diameter Accounting-Sub-Session-Id AVP.  This contains the accounting sub-session identifier.  The combination of the Session-Id and this AVP MUST be unique per sub-session, and the value of this AVP MUST be monotonically increased by one for all new sub-sessions.  The absence of this AVP implies no sub-sessions are in use, with the exception of an Accounting-Request whose Accounting-Record-Type is set to STOP_RECORD.  A STOP_RECORD message with no Accounting-Sub-Session-Id AVP present will signal the termination of all sub-sessions for a given Session-Id.</doc_description>
                <datatype>unsigned int64</datatype>
                <created_schema_version>4500</created_schema_version>
            </field>
            <field id='AccountingSessionId'>
                <doc_description>Diameter Acct-Session-Id AVP.  This is only used when RADIUS/Diameter translation occurs.  This contains the contents of the RADIUS Acct-Session-Id attribute.</doc_description>
                <datatype>string</datatype>
                <created_schema_version>4500</created_schema_version>
            </field>
            <field id='NormalizedSessionId'>
                <doc_description>This is a calculated value that can be used to identify a session.</doc_description>
                <datatype>unsigned int64</datatype>
                <created_schema_version>4500</created_schema_version>
            </field>
            <field id='SessionMsgId'>
                <doc_description>Diameter Accounting-Record-Number AVP.  This identifies this record within one session.  As Session-Id AVPs are globally unique, the combination of Session-Id and Accounting-Record-Number AVPs is also globally unique, and can be used in matching accounting records with confirmations.  An easy way to produce unique numbers is to set the value to 0 for records of type EVENT_RECORD and START_RECORD, and set the value to 1 for the first INTERIM_RECORD, 2 for the second, and so on until the value for STOP_RECORD is one more than for the last INTERIM_RECORD.</doc_description>
                <datatype>unsigned int32</datatype>
                <created_schema_version>4500</created_schema_version>
            </field>
            <field id='SourceHost'>
                <doc_description>Diameter Origin-Host AVP.  The Diameter Gateway knows that the Origin-Host AVP identifies the endpoint that originated the Diameter message.  However, the swapping of SourceHost/SourceRealm and DestinationHost/DestinationRealm fields is done in the Diameter Gateway.  Thus, outside the Diameter Gateway these fields do not have to be swapped.  Note that outside the Diameter Gateway, the SourceHost/SourceRealm fields represent the logical blade's side of the Diameter connection and the DestinationHost/DestinationRealm fields represent the other side of the Diameter connection.  If the SourceHost/SourceRealm fields are not present in the MDC when the Diameter Gateway processes it, the Diameter Gateway will add the appropriate configured values to the MDC.</doc_description>
                <datatype>string</datatype>
                <created_schema_version>4300</created_schema_version>
            </field>
            <field id='SourceRealm'>
                <doc_description>Diameter Origin-Realm AVP.  The Diameter Gateway knows that the Origin-Realm AVP identifies the realm of the endpoint that originated the Diameter message.  However, the swapping of SourceHost/SourceRealm and DestinationHost/DestinationRealm fields is done in the Diameter Gateway.  Thus, outside the Diameter Gateway these fields do not have to be swapped.  Note that outside the Diameter Gateway, the SourceHost/SourceRealm fields represent the logical blade's side of the Diameter connection and the DestinationHost/DestinationRealm fields represent the other side of the Diameter connection.  If the SourceHost/SourceRealm fields are not present in the MDC when the Diameter Gateway processes it, the Diameter Gateway will add the appropriate configured values to the MDC.</doc_description>
                <datatype>string</datatype>
                <created_schema_version>4300</created_schema_version>
            </field>
            <field id='DestinationRealm'>
                <doc_description>Diameter Destination-Realm AVP.  The Diameter Gateway knows that the Destination-Realm AVP identifies the realm of the destination of the Diameter message.  However, the swapping of SourceHost/SourceRealm and DestinationHost/DestinationRealm fields is done in the Diameter Gateway.  Thus, outside the Diameter Gateway these fields do not have to be swapped.  Note that outside the Diameter Gateway, the SourceHost/SourceRealm fields represent the logical blade's side of the Diameter connection and the DestinationHost/DestinationRealm fields represent the other side of the Diameter connection.  If the SourceHost/SourceRealm fields are not present in the MDC when the Diameter Gateway processes it, the Diameter Gateway will add the appropriate configured values to the MDC.</doc_description>
                <datatype>string</datatype>
                <created_schema_version>4300</created_schema_version>
            </field>
            <field id='DestinationHost'>
                <doc_description>Diameter Destination-Host AVP.  The Diameter Gateway knows that the Destination-Host AVP identifies the host of the destination of the Diameter message.  However, the swapping of SourceHost/SourceRealm and DestinationHost/DestinationRealm fields is done in the Diameter Gateway.  Thus, outside the Diameter Gateway these fields do not have to be swapped.  Note that outside the Diameter Gateway, the SourceHost/SourceRealm fields represent the logical blade's side of the Diameter connection and the DestinationHost/DestinationRealm fields represent the other side of the Diameter connection.  If the SourceHost/SourceRealm fields are not present in the MDC when the Diameter Gateway processes it, the Diameter Gateway will add the appropriate configured values to the MDC.</doc_description>
                <datatype>string</datatype>
                <created_schema_version>4300</created_schema_version>
            </field>
            <field id='SourceStateId'>
                <doc_description>Diameter Origin-State-Id AVP.  This is a monotonically increasing value that is advanced whenever a Diameter entity restarts with loss of previous state, for example upon reboot.  Origin-State-Id MAY be included in any Diameter message, including CER.  A Diameter entity issuing this AVP MUST create a higher value for this AVP each time its state is reset.  A Diameter entity MAY set Origin-State-Id to the time of startup, or it MAY use an incrementing counter retained in non-volatile memory across restarts.  The Origin-State-Id, if present, MUST reflect the state of the entity indicated by Origin-Host.  If a proxy modifies Origin-Host, it MUST either remove Origin-State-Id or modify it appropriately as well.  Typically, Origin-State-Id is used by an access device that always starts up with no active sessions; that is, any session active prior to restart will have been lost.  By including Origin-State-Id in a message, it allows other Diameter entities to infer that sessions associated with a lower Origin-State-Id are no longer active.  If an access device does not intend for such inferences to be made, it MUST either not include Origin-State-Id in any message, or set its value to 0.</doc_description>
                <datatype>unsigned int32</datatype>
                <created_schema_version>4500</created_schema_version>
            </field>
            <field id='AuthApplicationId'>
                <doc_description>Diameter Auth-Application-Id AVP.</doc_description>
                <datatype>unsigned int32</datatype>
                <created_schema_version>4300</created_schema_version>
            </field>
            <field id='ProxyInfo'>
                <doc_description>Diameter Proxy-Info AVP.</doc_description>
                <datatype>string</datatype>
                <created_schema_version>4300</created_schema_version>
                <deleted_schema_version>4610</deleted_schema_version>
            </field>
            <field id='RouteInfo'>
                <doc_description>Diameter Route-Record AVP.</doc_description>
                <datatype>string</datatype>
                <created_schema_version>4300</created_schema_version>
                <deleted_schema_version>4610</deleted_schema_version>
            </field>
            <field id='ExtraAvpList'>
                <doc_description>AVP.</doc_description>
                <datatype>blob</datatype>
                <created_schema_version>4300</created_schema_version>
                <list>1</list>
            </field>
            <field id='PricingDatabaseVersion'>
                <doc_description></doc_description>
                <datatype>unsigned int64</datatype>
                <created_schema_version>4300</created_schema_version>
                <deleted_schema_version>5122</deleted_schema_version>
            </field>
            <field id='BackstopFlag'>
                <doc_description>Indicates that this is a backstop message.</doc_description>
                <datatype>bool</datatype>
                <created_schema_version>4323</created_schema_version>
                <deleted_schema_version>4750</deleted_schema_version>
            </field>
            <field id='ProxyInfoList'>
                <doc_description>Diameter Proxy-Info AVP.</doc_description>
                <datatype>struct</datatype>
                <created_schema_version>4610</created_schema_version>
                <list>1</list>
                <struct_id>MtxDiamProxyInfoData</struct_id>
            </field>
            <field id='RouteInfoList'>
                <doc_description>Diameter Route-Record AVP.</doc_description>
                <datatype>string</datatype>
                <created_schema_version>4610</created_schema_version>
                <list>1</list>
            </field>
        </container>
    </configuration>
    `

	config, err := Parse(strings.NewReader(data))
	assert.Nil(t, err)

	assert.Equal(t, "Enum5GDeliveryReportRequested", config.Subtypes[0].ID)
	assert.Equal(t, "unsigned int16", config.Subtypes[0].Datatype)
	assert.Equal(t, "YES", config.Subtypes[0].Values[0].Value)
	assert.Equal(t, "NO", config.Subtypes[0].Values[1].Value)

	mtxDiamMsg := config.Containers[2]
	assert.Equal(t, "MtxDiamMsg", mtxDiamMsg.ID)

	assert.Equal(t, 94, mtxDiamMsg.Key)
	assert.Equal(t, 4300, mtxDiamMsg.CreatedSchemaVersion)
	assert.Equal(t, 0, mtxDiamMsg.DeletedSchemaVersion)
	assert.Equal(t, "4300", mtxDiamMsg.BaseContainer.ID)
	assert.Equal(t, "MtxChrgMsg", mtxDiamMsg.BaseContainer.Name)

	// First fields
	assert.Equal(t, "ProxyFlag", mtxDiamMsg.Fields[0].ID)
	assert.Equal(t, 4300, mtxDiamMsg.Fields[0].CreatedSchemaVersion)
	assert.Equal(t, 0, mtxDiamMsg.Fields[0].DeletedSchemaVersion)
	assert.Equal(t, "bool", mtxDiamMsg.Fields[0].Datatype)

	// Second fields
	assert.Equal(t, "ApplicationId", mtxDiamMsg.Fields[1].ID)
	assert.Equal(t, 4300, mtxDiamMsg.Fields[1].CreatedSchemaVersion)
	assert.Equal(t, 0, mtxDiamMsg.Fields[1].DeletedSchemaVersion)
	assert.Equal(t, "unsigned int32", mtxDiamMsg.Fields[1].Datatype)
}
