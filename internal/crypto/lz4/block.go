package lz4

// BlockUncompress uncompresses a lz4 block.
func BlockUncompress(src, dst []byte) (err error) {
	srcBuf, dstBuf := newBuffer(src), newBuffer(dst)
	var literalLen, matchLen, extraLen, offset int
	for err == nil {
		// Read token
		if literalLen, matchLen, err = srcBuf.ReadToken(); err != nil {
			break
		}
		// Optional: Copy literal to dst
		if literalLen > 0 {
			if literalLen == 0x0f {
				if extraLen, err = srcBuf.ReadExtraLength(); err != nil {
					break
				}
				literalLen += extraLen
			}
			if err = CopyN(dstBuf, srcBuf, literalLen); err != nil {
				break
			}
		}
		// Early return when reach the end
		if srcBuf.AtTheEnd() {
			break
		}
		// Read offset
		if offset, err = srcBuf.ReadOffset(); err != nil {
			break
		}
		// Update match length
		if matchLen == 0x13 {
			if extraLen, err = srcBuf.ReadExtraLength(); err != nil {
				break
			}
			matchLen += extraLen
		}
		dstBuf.WriteMatchedLiteral(offset, matchLen)
	}
	return
}
