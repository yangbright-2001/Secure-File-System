# Secure File Sharing System

A computer security project from [CS 161, Fall 2023](https://fa23.cs161.org/proj2/) at UC Berkeley, developed together with Chengle Zheng (team **Go Bears**).

## 1. What this project is

The course storage server is a shared cloud that anyone can read and rewrite. This repository is our client: it keeps files confidential and tamper-evident, and it supports sharing and revocation.

The client has to:

1. Log a user in with a password, and restore that user after the program exits.
2. Store, load, and append files without revealing the filename or the contents.
3. Share a file, and let the recipient choose their own filename.
4. Revoke one direct recipient and everyone that recipient passed the file to.

The behavior in this README is the design in [`fa23-proj2-go-bears-1/client/client.go`](fa23-proj2-go-bears-1/client/client.go). [`Design_Doc._v3.pdf`](Design_Doc._v3.pdf) describes the same structures. [`SUMMARY.docx`](SUMMARY.docx) is notes copied from the course spec.

Staff graded the client on Gradescope by running hidden cases. Our own cases are in [`fa23-proj2-go-bears-1/client_test/client_test.go`](fa23-proj2-go-bears-1/client_test/client_test.go). There is no separate input or output file. Section 3.2 copies one of those cases.

## 2. Repository layout

```
Design_Doc._v3.pdf                 design write-up
SUMMARY.docx                       spec notes
fa23-proj2-go-bears-1/
  client/client.go                 implementation
  client_test/client_test.go       examples
  README.md                        course starter README, including team SIDs
project2-userlib/                  course library: AES, RSA, HMAC, Argon2, Datastore, Keystore
```

Cryptography is called from `project2-userlib` (`SymEnc`, `PKEEnc`, `DSSign`, `HMACEval`, `Argon2Key`, `RandomBytes`). The client does not implement those primitives.

## 3. How to run

### 3.1 Tests

Go 1.20 or newer. From the test directory:

```bash
cd fa23-proj2-go-bears-1/client_test
go test -v
```

The first run downloads `github.com/cs161-staff/project2-userlib` and the test libraries. `go test` runs every case in `client_test/client_test.go`. Gradescope ran a larger hidden suite against the same API:

```text
InitUser(username, password) (*User, error)
GetUser(username, password) (*User, error)
User.StoreFile(filename, content) error
User.LoadFile(filename) ([]byte, error)
User.AppendToFile(filename, content) error
User.CreateInvitation(filename, recipientUsername) (invitationUUID, error)
User.AcceptInvitation(senderUsername, invitationUUID, filename) error
User.RevokeAccess(filename, recipientUsername) error
```

### 3.2 Example

These strings are the constants at the top of [`client_test/client_test.go`](fa23-proj2-go-bears-1/client_test/client_test.go). The case is `Testing Single User Store/Load/Append.`.

```text
InitUser("alice", "password")
StoreFile("aliceFile.txt", "Bitcoin is Nick's favorite ")
AppendToFile("aliceFile.txt", "digital ")
AppendToFile("aliceFile.txt", "cryptocurrency!")

LoadFile("aliceFile.txt")
  -> "Bitcoin is Nick's favorite digital cryptocurrency!"
```

The same file also checks sharing. Alice stores the file on one device, a second `GetUser("alice", "password")` shares it, and Bob accepts it as `bobFile.txt`. After Bob appends `"digital "` and Alice appends `"cryptocurrency!"`, Alice's two devices and Bob all load the same string above. Alice then revokes Bob. Alice can still load the file. Bob, and Charles if Bob had shared it onward, get an error on load and on append.

## 4. Data Structure & Design

Datastore is untrusted: any UUID can be read, changed, added, or deleted. Keystore holds public keys only, and an entry cannot be overwritten. Passwords have enough entropy for `Argon2Key`. A revoked user may have saved values from before revocation and may edit Datastore, and still must not read later updates.

### 4.1 How a record is sealed

Each secret value on Datastore is encrypt-then-MAC, or encrypt-then-sign. The tag is appended and has a fixed length, so checking integrity is a split.

```text
file node:   SymEnc(encKey, random IV, JSON) || HMAC(hmacKey, ciphertext)     HMAC is 64 bytes
user record: SymEnc(passwordKey, random IV, JSON) || signature                 signature is 256 bytes
invitation:  PKEEnc(recipient public key, JSON) || signature                   signature is 256 bytes
```

A file node's encryption key and HMAC key are two different 16-byte values from `RandomBytes`.

### 4.2 User

`InitUser` stores one `User` at `Hash(username)[:16]`. The struct holds the username, the password, an RSA decryption key (`Sharing_decrypt`), and a signing key (`Sharing_sign`).

```text
encryption key = Argon2Key(password, username + "salt", 16)
```

The username is the salt, so two users with the same password get different keys. Three public keys go to Keystore:

```text
username + "login"          verify key for the user record
username + "share_encrypt"  RSA key used to encrypt an invitation to this user
username + "share_sign"     verify key for invitations this user sends
```

`GetUser` verifies the signature, decrypts with the password, and checks the stored password. A second device calls `GetUser` again and gets its own `User`. Files live in Datastore, so that device sees updates immediately.

### 4.3 File layout

A user's filename is found by a deterministic address. The filename itself is not stored.

```text
UUID    = Hash(username + filename)[:16]
encKey  = Argon2Key(password, username + filename + "cryptography", 16)
hmacKey = Argon2Key(password, username + filename + "HMAC", 16)
```

That record is one of two types.

| | Owner, `File_1` | Recipient, `File_2` |
| --- | --- | --- |
| `type` | 1 | 2 |
| Points at | `File_head` | `Owner_share_metadata` (the shared channel) |
| Also stores | keys that open the header, plus a descendant list | keys that open the channel, plus the sender's username |

The file body is one `File_head` and a chain of `File_content` nodes. In the boxes below, the top row is the struct type. Every row under the line is one field of that struct. `"c"` is the value stored in the field `Content`, not part of the type name.

`File_head` points at the newest chunk. Each `File_content` points at the older one and stores the keys that open it. The oldest chunk, from the first `StoreFile`, has `Previous_uuid = nil`.

```text
+---------------------------------------------+
| type: File_1                                |
+---------------------------------------------+
| Head_ptr                         -> header  |
| Encryption_key                   header key |
| HMAC_key                         header key |
| File_type                        1          |
| DescendentList_ptr               -> list    |
| DescendentList_Encryption_key    list key   |
| DescendentList_HMAC_key          list key   |
+---------------------------------------------+
                       |
                       | field Head_ptr
                       v
+---------------------------------------------+
| type: File_head                             |
+---------------------------------------------+
| Content_ptr                      -> chunk   |
| Encryption_key                   key of "c" |
| HMAC_key                         key of "c" |
+---------------------------------------------+
                       |
                       | field Content_ptr
                       v
+---------------------------------------------+
| type: File_content                          |
+---------------------------------------------+
| Content                          "c"        |
| Previous_uuid                    -> older   |
| Encryption_key                   key of "b" |
| HMAC_key                         key of "b" |
+---------------------------------------------+
                       |
                       | field Previous_uuid
                       v
+---------------------------------------------+
| type: File_content                          |
+---------------------------------------------+
| Content                          "b"        |
| Previous_uuid                    -> older   |
| Encryption_key                   key of "a" |
| HMAC_key                         key of "a" |
+---------------------------------------------+
                       |
                       | field Previous_uuid
                       v
+---------------------------------------------+
| type: File_content                          |
+---------------------------------------------+
| Content                          "a"        |
| Previous_uuid                    nil        |
| Encryption_key                   nil        |
| HMAC_key                         nil        |
+---------------------------------------------+
```

There are three nodes of type `File_content`. They are the same struct. They differ only in field values: `Content` is `"c"`, then `"b"`, then `"a"`.

A recipient does not store a second copy of those chunks. `File_2.Share_ptr` points at one `Owner_share_metadata`, and that node's `Head_pointer` points at the same `File_head` the owner uses.

```text
+---------------------------------------------+
| type: File_2                                |
+---------------------------------------------+
| Share_ptr                        -> channel |
| Encryption_key                   channel key|
| HMAC_key                         channel key|
| File_type                        2          |
| Parent_name                      "alice"    |
+---------------------------------------------+
                       |
                       | field Share_ptr
                       v
+---------------------------------------------+
| type: Owner_share_metadata                  |
+---------------------------------------------+
| Head_pointer                     -> header  |
| Encryption_key                   header key |
| HMAC_key                         header key |
+---------------------------------------------+
                       |
                       | field Head_pointer
                       v
+---------------------------------------------+
| type: File_head                             |
+---------------------------------------------+
| Content_ptr                      -> chunk   |
| Encryption_key                   key of "c" |
| HMAC_key                         key of "c" |
+---------------------------------------------+
```

`LoadFile` follows those pointer fields, checks the HMAC, decrypts each node, and reverses the `Content` values. The caller receives `abc`.

### 4.4 Append and overwrite

`AppendToFile` rewrites the header and adds one chunk. It does not re-read the old chunks.

```text
1. Open the namespace record and the header.
2. Create a new chunk. Copy the header's current pointer and keys into it.
3. Point the header at the new chunk, and put fresh keys in the header.
4. Encrypt the new chunk with those fresh keys.
   Encrypt the header again with the header keys it already had.
```

`StoreFile` on a name that already exists checks the chain, then replaces it with one new chunk. The header UUID stays the same, so existing shares still find the file.

### 4.5 Sharing and revocation

Each person the owner shares with directly gets their own `Owner_share_metadata`. Its fields hold the header UUID and the keys that open the header. `Invitation_pointer` points at that node. The invitation blob is encrypted with the recipient's `share_encrypt` public key and signed with the sender's `Sharing_sign`.

```text
+---------------------------------------------+
| type: Invitation_pointer                    |
+---------------------------------------------+
| Meta                             -> channel |
| Encrypt                          channel key|
| HMAC                             channel key|
+---------------------------------------------+
                       |
                       | field Meta
                       v
+---------------------------------------------+
| type: Owner_share_metadata                  |
+---------------------------------------------+
| Head_pointer                     -> header  |
| Encryption_key                   header key |
| HMAC_key                         header key |
+---------------------------------------------+
                       |
                       | field Head_pointer
                       v
                 type: File_head
```

After `AcceptInvitation`, the recipient's `File_2` stores the same channel UUID in `Share_ptr`, plus the `Encrypt` and `HMAC` values that open it. The owner also records that recipient in the descendant list: channel UUID, those two keys, and the invitation UUID. The owner sends the invitation UUID outside this system.

A non-owner does not create a new `Owner_share_metadata`. Bob's invitation to Charles copies Bob's channel. Charles's `File_2.Share_ptr` and Bob's `File_2.Share_ptr` are the same UUID.

```text
File_head
|
|-- Owner_share_metadata          Bob's channel, keys K_bob
|     |-- File_2   username bob,    Share_ptr = this channel
|     '-- File_2   username charles, Share_ptr = this same channel
|
'-- Owner_share_metadata          Doris's channel, keys K_doris
      |-- File_2   username doris,  Share_ptr = this channel
      '-- File_2   username eve,    Share_ptr = this same channel
```

`RevokeAccess` is owner-only. Alice revokes Bob by drawing new header keys, re-encrypting the header at the same UUID, and writing those new keys into every channel that should survive (Doris). Bob's channel is left on the old keys, and Bob's invitation blob is deleted. HMAC verification then fails for Bob and for Charles. Doris and Eve still open the header. A second revoke of Bob returns an error.

Later appends use new chunk keys that exist only inside the re-encrypted header, so a revoked user who saved the old header keys cannot read those chunks.
