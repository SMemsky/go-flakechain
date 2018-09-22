# Levin Protocol docs

Note: this is just a quick writeup

Also look at Portable Storage docs: `pstorage.md`


## Ok, here we go

Each packet starts with a "bucket" or simply a `Header`.

Here it is, lol

### Packet types

* Notify
* Request
* Response

### Header

```
type bucketHead struct {
    Signature       uint64 // Should always be the right one :)
    PacketSize      uint64 // Exactly the size of the data following this header
    ReturnData      bool   // true for Request, false for others
    Command         uint32 // Command ID
    ReturnCode      int32  // Always zero for Notify and Request
    Flags           uint32 // 1 - Request, 2 - Response
    ProtocolVersion uint32 // Only version 1 is supported currently
}
```

By the way. Signature is always `0x0101010101012101`


Packet body (of size PacketSize) is a binary encoded struct (see `storages/portable`) and is written right after the header.


### Return code?

Not sure about it right not. Will check it later. TODO.


### Command ID

Packet identifier. Request and Response should have the same ID.