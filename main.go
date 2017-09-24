package main

import (
	"log"
	"net"
	"os"

	"github.com/gmule/gmule-core/protocol/ed2k"
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

const (
	// serverAddr = "ginuerzh.xyz:4661"
	// serverAddr = "176.103.56.98:2442"
	serverAddr = "176.103.48.36:4184"
)

func main() {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		log.Fatal(err)
	}

	var message ed2k.Message

	message = &ed2k.LoginMessage{
		UID: ed2k.UID([16]byte{0xd0, 0x97, 0x9a, 0x3d, 0x5c, 0x0e, 0xd6, 0x88, 0x75, 0xdc, 0x7d, 0x07, 0x8d, 0xc9, 0x6f, 0x20}),
		// ClientID: 0x212fa5b4,
		//UID:     protocol.NewUID(),
		Port: 4662,
		Tags: []ed2k.Tag{
			ed2k.StringTag(ed2k.TagName, "gmule", false),
			ed2k.Uint32Tag(ed2k.TagVersion, ed2k.EDonkeyVerion),
			ed2k.Uint32Tag(ed2k.TagServerFlags, ed2k.CapUnicode),
			ed2k.Uint32Tag(ed2k.TagEMuleVersion, ed2k.EMuleVersion),
		},
	}

	sendMessage(conn, message)

	for {
		m, err := ed2k.ReadMessage(conn, ed2k.CSTCPMessage)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(m)

		switch m.Type() {
		case ed2k.MessageIDChange:
			sendMessage(conn, &ed2k.GetServerListMessage{})
		case ed2k.MessageServerStatus:
			sendMessage(conn, offerFileMessage())
			sendMessage(conn, fileSearchMessage())
		}

	}
}

func offerFileMessage() (m ed2k.Message) {
	f, err := os.Open("test_file")
	if err != nil {
		return
	}

	hash, err := ed2k.Hash(f)
	if err != nil {
		return
	}
	file := ed2k.File{
		ClientID: 0xfbfbfbfb,
		Port:     0xfb,
	}
	copy(file.Hash[:], hash.Hash)
	file.Tags = []ed2k.Tag{
		ed2k.StringTag(ed2k.TagName|ed2k.TagCompactNameFlag, "test_file", true),
		ed2k.IntegerTag(ed2k.TagSize|ed2k.TagCompactNameFlag, uint64(hash.Size)),
	}
	m = &ed2k.OfferFilesMessage{
		Files: []ed2k.File{file},
	}
	return
}

func fileSearchMessage() ed2k.Message {
	return &ed2k.SearchRequestMessage{
		Searcher: ed2k.FileNameSearcher("mathematics"),
	}
}

func sendMessage(conn net.Conn, m ed2k.Message) (err error) {
	log.Println(m)
	data, err := m.Encode()
	if err != nil {
		return
	}
	if _, err = conn.Write(data); err != nil {
		return
	}
	return
}

/*
176.103.48.36:4184
0000   e3 4d 00 00 00 01 d0 97 9a 3d 5c 0e d6 88 75 dc  .M.......=\...u.
0010   7d 07 8d c9 6f 20 b4 a5 2f 21 36 12 04 00 00 00  }...o ../!6.....
0020   02 01 00 01 14 00 68 74 74 70 3a 2f 2f 77 77 77  ......http://www
0030   2e 61 4d 75 6c 65 2e 6f 72 67 03 01 00 11 3c 00  .aMule.org....<.
0040   00 00 03 01 00 20 1d 01 00 00 03 01 00 fb 00 10  ..... ..........
0050   04 03                                            ..
*/

/*
176.103.56.98:2442
0000   e3 4d 00 00 00 01 d0 97 9a 3d 5c 0e d6 88 75 dc  .M.......=\...u.
0010   7d 07 8d c9 6f 20 b4 a5 2f 21 36 12 04 00 00 00  }...o ../!6.....
0020   02 01 00 01 14 00 68 74 74 70 3a 2f 2f 77 77 77  ......http://www
0030   2e 61 4d 75 6c 65 2e 6f 72 67 03 01 00 11 3c 00  .aMule.org....<.
0040   00 00 03 01 00 20 1d 01 00 00 03 01 00 fb 00 10  ..... ..........
0050   04 03                                            ..

0000   e3 05 00 00 00 40 b4 a5 2f 21 e3 09 00 00 00 34  .....@../!.....4
0010   96 14 00 00 0b b6 19 00 e3 5a 00 00 00 38 57 00  .........Z...8W.
0020   73 65 72 76 65 72 20 76 65 72 73 69 6f 6e 20 31  server version 1
0030   37 2e 31 35 20 28 6c 75 67 64 75 6e 75 6d 29 0a  7.15 (lugdunum).
0040   50 6c 65 61 73 65 2c 20 6e 6f 74 65 20 74 68 61  Please, note tha
0050   74 20 65 44 6f 6e 6b 65 79 20 73 65 72 76 65 72  t eDonkey server
0060   73 20 64 6f 20 6e 6f 74 20 68 6f 73 74 20 61 6e  s do not host an
0070   79 20 66 69 6c 65 0a                             y file.
*/

/* search mathematics
0000   e3 0f 00 00 00 16 01 0b 00 6d 61 74 68 65 6d 61  .........mathema
0010   74 69 63 73                                      tics
*/

/*
0000   e3 15 00 00 00 19 28 72 aa bf 51 17 13 98 b2 c6  ......(r..Q.....
0010   1f 89 93 9e 62 61 4d 9c 4c 00                    ....baM.L.
*/

/* client callback request
0000   e3 05 00 00 00 1c 66 91 fd 00                    ......f...
*/
