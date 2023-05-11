package client

// CS 161 Project 2

// You MUST NOT change these default imports. ANY additional imports
// may break the autograder!

import (
	userlib "github.com/cs161-staff/project2-userlib"

	// Likewise, useful for debugging, etc.
	"encoding/hex"

	// Useful for string mainpulation.
	"strings"

	// Useful for formatting strings (e.g. `fmt.Sprintf`).
	"fmt"

	// Want to import errors.
	"errors"

	// Optional. You can remove the "_" there, but please do not touch
	// anything else within the import bracket.
	_ "strconv"
)

// This serves two purposes:
// a) It shows you some useful primitives, and
// b) it suppresses warnings for items not being imported.

func someUsefulThings() {
	// Creates a random UUID
	f := userlib.UUIDNew()
	userlib.DebugMsg("UUID as string:%v", f.String())

	// Example of writing over a byte of f
	f[0] = 10
	userlib.DebugMsg("UUID as string:%v", f.String())

	// Takes a sequence of bytes and renders as hex
	h := hex.EncodeToString([]byte("fubar"))
	userlib.DebugMsg("The hex: %v", h)

	// Marshals data into a JSON representation
	// Works well with Go structures!
	d, _ := userlib.Marshal(f)
	userlib.DebugMsg("The json data: %v", string(d))
	var g userlib.UUID
	userlib.Unmarshal(d, &g)
	userlib.DebugMsg("Unmashaled data %v", g.String())

	// errors.New(...) creates an error type!
	userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("This is an error")))

	// Generate a random private/public keypair.
	// The "_" indicates that we don't check for the error case here.
	var pk userlib.PKEEncKey
	var sk userlib.PKEDecKey
	pk, sk, _ = userlib.PKEKeyGen()
	userlib.DebugMsg("Key is %v, %v", pk, sk)

	// Useful for string interpolation.
	_ = fmt.Sprintf("%s_%d", "file", 1)
}

var KEY_LENGTH = 16
var NULL_UUID, _ = userlib.UUIDFromBytes([]byte{0})

func IsSliceEqual(a, b []byte) bool {
	for i, n := range b {
		if a[i] != n {
			return false
		}
	}
	return true
}

type Node struct {
	element []byte
	next *Node
}

type Queue struct {
	first *Node
	last *Node
}

func (q *Queue) Peek() []byte { return q.first.element }

func (q *Queue) IsEmpty() bool { return q.first==nil }

func (q *Queue) Clear() {
	for q.first == nil {
		// destroy first
		q.first = q.first.next
	}
	q.last = nil
}

func (q *Queue) Push(element []byte) {
	var node Node
	node.element = element
	node.next = nil
	if (q.first == nil) {
		q.first = &node
	} else {
		q.last.next = &node
	}
	q.last = &node
}

func (q *Queue) Pop() (element []byte) {
	if (q.IsEmpty()) { return nil }
	element = q.first.element
	q.first = q.first.next
	return element
}

type Integer struct {
	I int
}

func MarshalInt(i int) (arr []byte, err error) {
	var wrapper Integer
	wrapper.I = i
	return userlib.Marshal(wrapper)
}
func UnmarshalInt(arr []byte) (i int, err error) {
	var wrapper Integer
	err = userlib.Unmarshal(arr, &wrapper)
	return wrapper.I, err
}

type FileMetadata struct {
	Owned bool
	FileKey []byte
}

// User is the structure definition for a user record.
type User struct {
	Username string
	KUser []byte // the cryptographic key generated from the username and password
	DecryptKey userlib.PKEDecKey // for public key encryption
	SignKey userlib.DSSignKey // for digital signatures
}

type DatastoreEntry struct {
	Ciphertext []byte
	VerifyBytes []byte
}

func GetHmac(key []byte, value []byte) (hmac []byte, err error) {
	hmacKey, err := userlib.HashKDF(key[:KEY_LENGTH], []byte("HMAC"))
	if (err != nil) {return nil, err}
	return userlib.HMACEval(hmacKey[:KEY_LENGTH], value)
}

func GetEntryHmac(key, value []byte) (entry DatastoreEntry, err error) {
	entry.Ciphertext = userlib.SymEnc(key[:KEY_LENGTH], userlib.RandomBytes(KEY_LENGTH), value)
	entry.VerifyBytes, err = GetHmac(key[:KEY_LENGTH], entry.Ciphertext)
	return entry, err
}

func GetEntrySigned(encryptKey userlib.PKEEncKey, signedKey userlib.DSSignKey, value []byte) (entry DatastoreEntry, err error) {
	entry.Ciphertext, err = userlib.PKEEnc(encryptKey, value)
	if (err != nil) {return entry, err}
	entry.VerifyBytes, err = userlib.DSSign(signedKey, entry.Ciphertext)
	return entry, err
}

func DataStoreEncryptSignedSet(encryptKey userlib.PKEEncKey, signedKey userlib.DSSignKey, uuid userlib.UUID, value []byte) (err error){
	var entry DatastoreEntry
	entry, err = GetEntrySigned(encryptKey, signedKey, value)
	if (err != nil) { return err }
	entryMarshaled, err := userlib.Marshal(entry)
	if (err != nil) { return err }
	
	userlib.DatastoreSet(uuid, entryMarshaled)
	return nil
}

func DataStoreDecryptSignedGet(decryptKey userlib.PKEDecKey, verifyKey userlib.DSVerifyKey, uuid userlib.UUID) (value []byte, err error) {
	jsonData, ok := userlib.DatastoreGet(uuid)
	if (!ok) { return nil, errors.New(strings.ToTitle("There is no data at specified uuid"))}
	
	var entry DatastoreEntry
	err = userlib.Unmarshal(jsonData, &entry)
	if (err != nil) { return nil, err }
	
	err = userlib.DSVerify(verifyKey, entry.Ciphertext, entry.VerifyBytes)
	if (err != nil) { return nil, err }

	valueDatastore, err := userlib.PKEDec(decryptKey, entry.Ciphertext)
	if (err != nil) { return nil, err }
	return valueDatastore, nil
}

func DataStoreEncryptSet(key []byte, value []byte) (err error) {
	var entry DatastoreEntry
	entry, err = GetEntryHmac(key, value)
	if (err != nil) { return err }
	
	address, err := userlib.UUIDFromBytes(userlib.Hash(key))
	if (err != nil) { return err }
	
	entryMarshaled, err := userlib.Marshal(entry)
	if (err != nil) { return err }
	
	userlib.DatastoreSet(address, entryMarshaled)
	return nil
}

func DataStoreDecryptGet(key []byte) (value []byte, err error) {
	var entry DatastoreEntry
	address, err := userlib.UUIDFromBytes(userlib.Hash(key))
	if (err != nil) { return nil, err }
	
	jsonData, ok := userlib.DatastoreGet(address)
	if (!ok) { return nil, errors.New(strings.ToTitle("No data at address")) }
	
	err = userlib.Unmarshal(jsonData, &entry)
	if (err != nil) { return nil, err }
	
	// check hmac first
	hmac, err := GetHmac(key, entry.Ciphertext)
	if (err != nil) {return nil, err}
	
	if (!IsSliceEqual(entry.VerifyBytes, hmac)) {
		return nil, errors.New(strings.ToTitle("Encrypted value has been changed maliciously: HMAC did not match"))
	}
	
	value = userlib.SymDec(key[:KEY_LENGTH], entry.Ciphertext)
	return value, nil
}

func DataStoreDeleteSet(key []byte) (err error) {
	address, err := userlib.UUIDFromBytes(userlib.Hash(key))
	if (err != nil) { return err }
	userlib.DatastoreDelete(address)
	return nil
}

func VerifyUsername (username string) (verifyKey userlib.DSVerifyKey, err error) {
	verifyKey, ok := userlib.KeystoreGet(KeyGenDSPublicKey(username))
	if (!ok) {return verifyKey, errors.New(strings.ToTitle("User does not exist!"))}
	return verifyKey, nil
}

func KeyGenPublicKey (username string) (string) {
	return string(userlib.Hash([]byte(username + "encrypt"))[:KEY_LENGTH])
}
func KeyGenDSPublicKey (username string) (string) {
	return string(userlib.Hash([]byte(username + "signatures"))[:KEY_LENGTH])
}

func GetFile(userdata User, filename string) (file FileMetadata, err error) {
	key, err := userlib.HashKDF(userdata.KUser[:KEY_LENGTH], append([]byte(filename), []byte("file")...))
	if (err != nil) { return file, err }
	jsonData, err := DataStoreDecryptGet(key)
	if (err != nil) { return file, err }
	err = userlib.Unmarshal(jsonData, &file)
	if (err != nil) { return file, err }
	
	if (!file.Owned) {
		_, err = DataStoreDecryptGet(file.FileKey)
		if (err != nil) {
			DataStoreDeleteSet(key)
			return file, errors.New(strings.ToTitle("Access to this file was revoked"))
		}
	}
	
	return file, nil
}

func SetFile(userdata User, filename string, file FileMetadata) (err error) {
	jsonData, err := userlib.Marshal(file)
	if (err != nil) { return err }
	key, err := userlib.HashKDF(userdata.KUser[:KEY_LENGTH], append([]byte(filename), []byte("file")...))
	if (err != nil) { return err }

	err = DataStoreEncryptSet(key, jsonData)
	return err
}

func GenKFile(userdata User, filename string) ([]byte, error) {
	return userlib.HashKDF(userdata.KUser[:KEY_LENGTH], append([]byte(filename), []byte("kFile")...))
}

func GetKFile(userdata User, filename string) (kFile []byte, err error) {
	fileMetadata, err := GetFile(userdata, filename)
	if (err != nil) { return nil, err }
	if (fileMetadata.Owned) {
		return GenKFile(userdata, filename)
	}
	
	kFile, err = DataStoreDecryptGet(fileMetadata.FileKey)
	return kFile, err
}

func GenInvitationLink(filePointer []byte) (uuid userlib.UUID, err error) {
	tempHash, err := userlib.HashKDF(filePointer[:KEY_LENGTH], []byte{1})
	address, err := userlib.UUIDFromBytes(tempHash)
	return address, err
}

func GenFilePointer(userdata User, filename string, username string) (result []byte, err error) {
	kFile, err := GetKFile(userdata, filename)
	if (err != nil) {return nil, err}
	
	result, err = userlib.HashKDF(kFile[:KEY_LENGTH], []byte(username))
	
	return result, err
}

func GetFilePointer(userdata User, filename string) (filePointer []byte, err error) {
	fileMetadata, err := GetFile(userdata, filename)
	if (err != nil) {return nil, err}
	if (fileMetadata.Owned) {
		return GenKFile(userdata, filename)
	}
	return fileMetadata.FileKey, nil
}

func GetKUser(username, password string) []byte {
	byteUsername := []byte(username)
	bytePassword := []byte(password)
	byteUser := append(byteUsername, bytePassword...)
	return userlib.Argon2Key(byteUser, []byte{1, 2, 3, 4, 5, 6}, uint32(KEY_LENGTH))
}

func InitUser(username string, password string) (userdataptr *User, err error) {
	// check username and password is fine
	if (len(username) == 0) {
		return nil, errors.New(strings.ToTitle("Username must contain at least one character"))
	}
	
	_, ok := userlib.KeystoreGet(KeyGenPublicKey(username))
	
	if (ok) {
		return nil, errors.New(strings.ToTitle("Username taken: All username must be unique"))
	}
	
	var userdata User
	userdata.Username = username
	
	var publicKey userlib.PKEEncKey
	var privateKey userlib.PKEDecKey
	publicKey, privateKey, err = userlib.PKEKeyGen()
	if (err != nil) { return nil, err }
	
	var DSPublicKey userlib.DSVerifyKey
	var DSPrivateKey userlib.DSSignKey
	DSPrivateKey, DSPublicKey, err = userlib.DSKeyGen()
	if (err != nil) { return nil, err }
	
	userlib.KeystoreSet(KeyGenPublicKey(username), publicKey)
	userlib.KeystoreSet(KeyGenDSPublicKey(username), DSPublicKey)

	userdata.DecryptKey = privateKey
	userdata.SignKey = DSPrivateKey
	
	byteUsername := []byte(username)
	bytePassword := []byte(password)
	byteUser := append(byteUsername, bytePassword...)
	userdata.KUser = GetKUser(username, password)
	
	// UUID(username||password||1) -> UUID(username||password||2)
	hash1 := userlib.Hash(append(byteUser, []byte{1}...))
	hash2 := userlib.Hash(append(byteUser, []byte{2}...))
	
	// key to encrypt: hash1
	// userlib.UUIDFromBytes(Hash(hash1)) as key to datastore
	err = DataStoreEncryptSet(hash1, hash2)
	if (err != nil) { return nil, err }
	
	userdataMarshaled, err := userlib.Marshal(userdata)
	if (err != nil) { return nil, err }
	
	err = DataStoreEncryptSet(userdata.KUser, userdataMarshaled)
	if (err != nil) { return nil, err }
	
	return &userdata, nil
}

func GetUser(username string, password string) (userdataptr *User, err error) {
	_, err = VerifyUsername(username)
	if (err != nil) { return nil, errors.New(strings.ToTitle("Username and passwrd did not match")) }
	
	var userdata User
	userdataptr = &userdata

	byteUsername := []byte(username)
	bytePassword := []byte(password)
	byteUser := append(byteUsername, bytePassword...)
	hash1 := userlib.Hash(append(byteUser, []byte{1}...))
	hash2 := userlib.Hash(append(byteUser, []byte{2}...))
	hash2Temp, err := DataStoreDecryptGet(hash1)
	if (err != nil) { return nil, err }
	
	if (!IsSliceEqual(hash2, hash2Temp)) {
		return nil, errors.New(strings.ToTitle("Username and password did not match"))
	}
	
	kUser := GetKUser(username, password)
	
	userdata_temp, err := DataStoreDecryptGet(kUser)
	if (err != nil) { return nil, err }
	
	err = userlib.Unmarshal(userdata_temp, userdataptr)
	if (err != nil) { return nil, err }
	
	return userdataptr, nil
}

func (userdata *User) StoreFile(filename string, content []byte) (err error) {
	kFile, err := GenKFile(*userdata, filename)
	if (err != nil) { return err }

	metadata, exist := GetFile(*userdata, filename)
	if (exist != nil) {
		metadata.Owned = true
		metadata.FileKey = nil
		err = SetFile(*userdata, filename, metadata)
		if (err != nil) { return err }
		err = SetSharedList(kFile, [][]byte{})
		if (err != nil) { return err }
	}

	one, err := MarshalInt(1)
	if (err != nil) { return err }
	err = DataStoreEncryptSet(kFile, one)
	if (err != nil) { return err }

	zero, err := MarshalInt(0)
	if (err != nil) { return err }
	dataKey, err := userlib.HashKDF(kFile[:KEY_LENGTH], zero)
	if (err != nil) { return err }
	
	err = DataStoreEncryptSet(dataKey, content)
	if (err != nil) { return err }
	
	return nil
}

func (userdata *User) AppendToFile(filename string, content []byte) (err error) {
	storageKey, err := GetKFile(*userdata, filename)
	if (err != nil) { return err }
	
	length, err := DataStoreDecryptGet(storageKey)
	if (err != nil) { return err }
	
	var lengthInt int
	lengthInt, err = UnmarshalInt(length)
	if (err != nil) { return err }
	appendKey, err := userlib.HashKDF(storageKey[:KEY_LENGTH], length)
	if (err != nil) { return err }
	err = DataStoreEncryptSet(appendKey, content)
	if (err != nil) { return err }
	
	lengthInt++
	length, err = MarshalInt(lengthInt)
	if (err != nil) { return err }
	err = DataStoreEncryptSet(storageKey, length)
	if (err != nil) { return err }
	
	return nil
}

func (userdata *User) LoadFile(filename string) (content []byte, err error) {
	storageKey, err := GetKFile(*userdata, filename)
	if (err != nil) { return nil, err }
	
	length, err := DataStoreDecryptGet(storageKey)
	if (err != nil) { return nil, err }
	
	lengthInt, err := UnmarshalInt(length)
	if (err != nil) { return nil, err }
	
	for i := 0; i < lengthInt; i++ {
		iMarshal, err := MarshalInt(i)
		if (err != nil) { return nil, err }
		
		dataKey, err := userlib.HashKDF(storageKey[:KEY_LENGTH], iMarshal)
		if (err != nil) { return nil, err }

		contentPart, err := DataStoreDecryptGet(dataKey)
		if (err != nil) { return nil, err }
		
		content = append(content, contentPart...)
	}
	
	if (len(content) == 0) {
		return make([]byte, 0), nil
	}
	
	return content, nil
}

func GetSharedListKey(filePointer []byte) ([]byte, error) {
	return userlib.HashKDF(filePointer[:KEY_LENGTH], []byte("Share List"))
}

func GetSharedList(filePointer []byte) (list [][]byte, err error) {
	sharedListKey, err := GetSharedListKey(filePointer)
	if (err != nil) { return nil, err}
	
	jsonData, err := DataStoreDecryptGet(sharedListKey)
	if (err != nil) { return nil, err }
	
	err = userlib.Unmarshal(jsonData, &list)
	if (err != nil) { return nil, err }
	return list, nil
}

func SetSharedList(filePointer []byte, sharedList [][]byte) (err error) {
	sharedListKey, err := GetSharedListKey(filePointer)
	if (err != nil) { return err }
	
	jsonData, err := userlib.Marshal(sharedList)
	if (err != nil) { return err }
	
	err = DataStoreEncryptSet( sharedListKey, jsonData )
	if (err != nil) {return err }
	
	return nil
}

func (userdata *User) CreateInvitation(filename string, recipientUsername string) (
	invitationPtr userlib.UUID, err error) {
	
	filePointer, err := GetFilePointer(*userdata, filename)
	if (err != nil) {return NULL_UUID, err}
	
	filePointerRecipient, err := GenFilePointer(*userdata, filename, recipientUsername)
	if (err != nil) {return NULL_UUID, err}
	
	invLink, err := GenInvitationLink(filePointerRecipient)
	if (err != nil) {return NULL_UUID, err}
	
	encryptKey, ok := userlib.KeystoreGet(KeyGenPublicKey(recipientUsername))
	if (!ok) { return NULL_UUID, errors.New(strings.ToTitle("Pubic key does not exist")) }
	
	kFile, err := GetKFile(*userdata, filename)
	if (err != nil) {return NULL_UUID, err}
	
	err = DataStoreEncryptSignedSet(encryptKey, userdata.SignKey, invLink, filePointerRecipient)
	if (err != nil) {return NULL_UUID, err}
	
	err = DataStoreEncryptSet(filePointerRecipient, kFile)
	if (err != nil) {return NULL_UUID, err}
	
	// add recipientUsername to below:
	sharedList, err := GetSharedList(filePointer)
	if (err != nil) { return NULL_UUID, err }
	
	sharedList = append(sharedList, filePointerRecipient)
	err = SetSharedList(filePointer, sharedList)
	if (err != nil) { return NULL_UUID, err }
	
	return invLink, err
}

func (userdata *User) AcceptInvitation(senderUsername string, invitationPtr userlib.UUID, filename string) (err error) {
	metadata, err := GetFile(*userdata, filename)
	if (err == nil) { return errors.New(strings.ToTitle("Filename already exists"))}
	
	verifyKey, err := VerifyUsername(senderUsername)
	if (err != nil) {return err}

	filePointer, err := DataStoreDecryptSignedGet(userdata.DecryptKey, verifyKey, invitationPtr)
	if (err != nil) {return err}
	
	metadata.Owned = false
	metadata.FileKey = filePointer
	err = SetFile(*userdata, filename, metadata)
	if (err != nil) {return err}
	
	err = SetSharedList(filePointer, [][]byte{})
	if (err != nil) { return err }
	
	return nil
}

func (userdata *User) RevokeAccess(filename string, recipientUsername string) (err error) {
	filePointer, err := GenFilePointer(*userdata, filename, recipientUsername)
	if (err != nil) { return err }
	
	
	var queue Queue
	queue.Push(filePointer)
	
	for (!queue.IsEmpty()) {
		//poping queue gives filePointer
		filePointer := queue.Pop()
		
		// get all other links that were sent by this guy, and add it to queue
		var sharedList [][]byte
		sharedList, err = GetSharedList(filePointer)
		// If below is disabled, because if file is shared but is not accepted,
		// sharedList might not exist
		// if (err != nil) {return err}
		// for _, n := range sharedList {
		// 	queue.Push(n)
		// }
		
		if (err == nil) {
			for _, n := range sharedList {
				queue.Push(n)
			}
		}
		
		sharedListKey, err := GetSharedListKey(filePointer)
		if (err != nil) { return err }
		_ = DataStoreDeleteSet(sharedListKey)
		
		// revoke access to curr
		err = DataStoreDeleteSet(filePointer)
		if (err != nil) { return err }
		
		invLink, err := GenInvitationLink(filePointer)
		if (err != nil) { return err }
		
		userlib.DatastoreDelete(invLink)
		
	}
	
	return nil
}
