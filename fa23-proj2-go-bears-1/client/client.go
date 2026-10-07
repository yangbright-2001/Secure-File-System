package client

// CS 161 Project 2

// Only the following imports are allowed! ANY additional imports
// may break the autograder!
// - bytes
// - encoding/hex
// - encoding/json
// - errors
// - fmt
// - github.com/cs161-staff/project2-userlib
// - github.com/google/uuid
// - strconv
// - strings

import (
	"encoding/json"

	userlib "github.com/cs161-staff/project2-userlib"
	"github.com/google/uuid"

	// hex.EncodeToString(...) is useful for converting []byte to string

	// Useful for string manipulation
	//"strings"

	// Useful for formatting strings (e.g. `fmt.Sprintf`).
	"fmt"

	// Useful for creating new error messages to return using errors.New("...")
	"errors"

	// Optional.
	_ "strconv"
)

// This serves two purposes: it shows you a few useful primitives,
// and suppresses warnings for imports not being used. It can be
// safely deleted!
func someUsefulThings() {

	// Creates a random UUID.
	randomUUID := uuid.New()

	// Prints the UUID as a string. %v prints the value in a default format.
	// See https://pkg.go.dev/fmt#hdr-Printing for all Golang format string flags.
	userlib.DebugMsg("Random UUID: %v", randomUUID.String())

	// Creates a UUID deterministically, from a sequence of bytes.
	hash := userlib.Hash([]byte("user-structs/alice"))
	deterministicUUID, err := uuid.FromBytes(hash[:16])
	if err != nil {
		// Normally, we would `return err` here. But, since this function doesn't return anything,
		// we can just panic to terminate execution. ALWAYS, ALWAYS, ALWAYS check for errors! Your
		// code should have hundreds of "if err != nil { return err }" statements by the end of this
		// project. You probably want to avoid using panic statements in your own code.
		panic(errors.New("An error occurred while generating a UUID: " + err.Error()))
	}
	userlib.DebugMsg("Deterministic UUID: %v", deterministicUUID.String())

	// Declares a Course struct type, creates an instance of it, and marshals it into JSON.
	type Course struct {
		name      string
		professor []byte
	}

	course := Course{"CS 161", []byte("Nicholas Weaver")}
	courseBytes, err := json.Marshal(course)
	if err != nil {
		panic(err)
	}

	userlib.DebugMsg("Struct: %v", course)
	userlib.DebugMsg("JSON Data: %v", courseBytes)

	// Generate a random private/public keypair.
	// The "_" indicates that we don't check for the error case here.
	var pk userlib.PKEEncKey
	var sk userlib.PKEDecKey
	pk, sk, _ = userlib.PKEKeyGen()
	userlib.DebugMsg("PKE Key Pair: (%v, %v)", pk, sk)

	// Here's an example of how to use HBKDF to generate a new key from an input key.
	// Tip: generate a new key everywhere you possibly can! It's easier to generate new keys on the fly
	// instead of trying to think about all of the ways a key reuse attack could be performed. It's also easier to
	// store one key and derive multiple keys from that one key, rather than
	originalKey := userlib.RandomBytes(16)
	derivedKey, err := userlib.HashKDF(originalKey, []byte("mac-key"))
	if err != nil {
		panic(err)
	}
	userlib.DebugMsg("Original Key: %v", originalKey)
	userlib.DebugMsg("Derived Key: %v", derivedKey)

	// A couple of tips on converting between string and []byte:
	// To convert from string to []byte, use []byte("some-string-here")
	// To convert from []byte to string for debugging, use fmt.Sprintf("hello world: %s", some_byte_arr).
	// To convert from []byte to string for use in a hashmap, use hex.EncodeToString(some_byte_arr).
	// When frequently converting between []byte and string, just marshal and unmarshal the data.
	//
	// Read more: https://go.dev/blog/strings

	// Here's an example of string interpolation!
	_ = fmt.Sprintf("%s_%d", "file", 1)
}

// This is the type definition for the User struct.
// A Go struct is like a Python or Java class - it can have attributes
// (e.g. like the Username attribute) and methods (e.g. like the StoreFile method below).

type User struct {
	Username        string
	Password        string
	Sharing_decrypt userlib.PKEDecKey //private key for encryption(生成一个key pair，private的在这，public的在keystore)
	Sharing_sign    userlib.DSSignKey //private key for signature

	// You can add other attributes here if you want! But note that in order for attributes to
	// be included when this struct is serialized to/from JSON, they must be capitalized.

	// On the flipside, if you have an attribute that you want to be able to access from
	// this struct's methods, but you DON'T want that value to be included in the serialized value
	// of this struct that's stored in datastore, then you can use a "private" variable (e.g. one that
	// begins with a lowercase letter).
}

// 记录单个descendent的数据结构
type Descendent struct {
	// Child_name 							string
	Child_metadata_ptr            uuid.UUID
	Child_metadata_Encryption_key []byte //用来解密metadata的
	Child_metadata_HMAC_key       []byte
	Invitationptr_uuid            uuid.UUID	//用来删库
}

// 记录所有descendents的数据结构
type Descendent_list struct {
	All_descendents map[string]Descendent //{child_name:descendent_struct}
}

type File_1 struct { //metadata  of owner
	Head_ptr                      uuid.UUID ///////////////////////本来是*File_head
	Encryption_key                []byte    //random byte generator
	HMAC_key                      []byte    //random byte generator
	File_type                     int
	DescendentList_ptr            uuid.UUID
	DescendentList_Encryption_key []byte
	DescendentList_HMAC_key       []byte
}

type File_2 struct { //metadata of the shared persons
	Share_ptr      uuid.UUID /////////////////////元素的
	Encryption_key []byte    //random byte generator
	HMAC_key       []byte
	File_type      int
	Parent_name    string //只要parent
}

type File_head struct {
	Content_ptr    uuid.UUID
	Encryption_key []byte //generate by RBG
	HMAC_key       []byte
}

// //////////把invitation_pointer和metadata拆成两个struct算了，不然混乱
type Owner_share_metadata struct { //A create invitation的时候，会产生一个含有head uuid的file metadata的copy
	Head_pointer   uuid.UUID //A给B就是指向head的uuid
	Encryption_key []byte    //generate by RBG
	HMAC_key       []byte
}

type Invitation_pointer struct { //invitation_pointer，用同一个structure
	Meta    uuid.UUID //指向那个metada的uuid
	Encrypt []byte    //generate by RBG
	HMAC    []byte
}

type File_content struct {
	Content        string
	Previous_uuid  uuid.UUID
	Encryption_key []byte //generate by RBG
	HMAC_key       []byte
}

// NOTE: The following methods have toy (insecure!) implementations.

func info_string2uuid(info_string string) (uuid.UUID, error) {
	info_byte := []byte(info_string)
	info_hash := userlib.Hash(info_byte)
	info_uuid, err_bytes := uuid.FromBytes(info_hash[:16]) //不安全，转为16byte的uuid
	if err_bytes != nil {
		err := errors.New("in info_string2uuid(): generate uuid failed")
		return uuid.Nil, err
	}
	return info_uuid, nil
}

// 把各个node中的content（每个node的content都是[]byte)按node翻转，然后正确的输出来（flatten二维切片）
func reverse_flatten(s [][]byte) []byte {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	s_flatten := make([]byte, 0)
	for _, v := range s {
		s_flatten = append(s_flatten, v...)
	}
	return s_flatten
}

// 利用password和salt加密user struct
func encrypt_userstruct_salt_signature(login_signature_key userlib.DSSignKey, password string, username string, userdata User) (encrypted_userdata_final []byte, err error) {
	///////////////////////////////marshal写法，参数是写userdata还是&userdata
	userdata_byte, err_marshal := json.Marshal(&userdata)
	if err_marshal != nil {
		err = errors.New("in encrypt_userstruct_salt_signature(): marshal serilization failed")
		return nil, err
	}

	//随机生成一个iv  //hash之后的password作为key
	//通过修改salt （username->username++） 和同样的password来形成一个新的symmetric key，用来加密user struct
	password_byte := []byte(password)
	salt_encrypt := []byte(username + "salt") //salt保证同一个密码的两个人不会same hash
	encrypt_symmetric_key := userlib.Argon2Key(password_byte, salt_encrypt, 16)
	IV := userlib.RandomBytes(16)
	encrypted_userdata := userlib.SymEnc(encrypt_symmetric_key, IV, userdata_byte)

	//做signature，连在加密的user struct后面
	signature_login, err_signature := userlib.DSSign(login_signature_key, encrypted_userdata)
	if err_signature != nil {
		err = errors.New("in encrypt_userstruct_salt_signature(): generate signature failed")
		return nil, err
	}
	encrypted_userdata_final = append(encrypted_userdata, signature_login...)
	return encrypted_userdata_final, nil
}

// //////////////////用deterministic的symmetric key加密filestruct(因为是symmetric，所以不用存)
// ///////////////////////////////////////////////////不论filestruct1和filestruct2都是用的同一种方法加密？
// 随机生成一个iv  //hash之后的password作为key
// 通过修改salt （username+filename+"__"） 和同样的password来形成一个新的symmetric key，用来加密file struct
// 确定好salt是统一的
func encrypt_filestruct_salt_HMAC(password string, username string, filename string, file_struct interface{}) (encrypted_filestruct_final []byte, err error) {
	filestruct_byte, err_marshal := json.Marshal(&file_struct)
	if err_marshal != nil {
		err = errors.New("in encrypt_filestruct_salt_HMAC(): marshal serilization for file_struct failed")
		return nil, err
	}
	password_byte := []byte(password)
	salt_encrypt := []byte(username + filename + "cryptography") //salt保证同一个密码的两个人不会same hash
	salt_HMAC := []byte(username + filename + "HMAC")
	encrypt_symmetric_key := userlib.Argon2Key(password_byte, salt_encrypt, 16)
	HMAC_key := userlib.Argon2Key(password_byte, salt_HMAC, 16)
	IV := userlib.RandomBytes(16)
	encrypted_filestruct := userlib.SymEnc(encrypt_symmetric_key, IV, filestruct_byte) //////////////
	HMAC_filestruct, err_HMAC := userlib.HMACEval(HMAC_key, encrypted_filestruct)
	if err_HMAC != nil {
		err = errors.New("in encrypt_filestruct_salt_HMAC(): the HMAC generation for file_struct failed")
		return nil, err
	}
	encrypted_filestruct_final = append(encrypted_filestruct, HMAC_filestruct...)
	return encrypted_filestruct_final, nil
}

// 对file node用symmetric encryption key（iv）加密+HMAC签名
// //////////////////////////////////如果interface{}参数不行的话，就写两个函数吧，head和content的加密各一个
func encrypt_HMAC_filenode(filenode interface{}, encryption_key []byte, HMAC_key []byte) ([]byte, error) {
	//序列化filenode
	var err error
	filenode_byte, err_marshal := json.Marshal(&filenode)
	if err_marshal != nil {
		err = errors.New("in encrypt_HMAC_filenode(): marshal serilization failed")
		return nil, err //切片的零值是nil
	}
	//帮filenode进行symmetric encryption
	IV := userlib.RandomBytes(16)
	encrypted_filenode := userlib.SymEnc(encryption_key, IV, filenode_byte)
	//生成HMAC
	HMAC_filenode, err_HMAC_filenode := userlib.HMACEval(HMAC_key, encrypted_filenode)
	if err_HMAC_filenode != nil {
		err = errors.New("in encrypt_HMAC_filenode(): generate HMAC failed")
		return nil, err
	}
	encrypted_filenode_final := append(encrypted_filenode, HMAC_filenode...)
	return encrypted_filenode_final, nil
}

// 用receiver的encryption key去加密，用caller的sign key去sign，对象暂时确定是invitation pointer
func encrypt_signature_invitationptr(invitation_pointer Invitation_pointer, recipient_encryption_key userlib.PKEEncKey, caller_sign_key userlib.DSSignKey, recipientUsername string) (encrypted_invitationptr_final []byte, err error) {
	invitationptr_byte, err_marshal := json.Marshal(invitation_pointer)
	if err_marshal != nil {
		err = errors.New("in encrypt_signature_invitationptr(): marshal serilization for invitation_pointer failed")
		return nil, err
	}

	encrypted_invitationptr, err_encryption := userlib.PKEEnc(recipient_encryption_key, invitationptr_byte)
	if err_encryption != nil {
		err = errors.New("in encrypt_signature_invitationptr(): public_key encryption for invitation_pointer failed")
		return nil, err
	}

	/////////////////////////////////////////////////////////////////////////////////
	//////////	直接pKEEnc有时候会报错，这里考虑改成hybrid encryption	////////////////////////////
	/////////////////////////////////////////////////////////////////////////////////
	// ////hybrid encryption
	// //generate a random symmetric key
	// random_sym_key := userlib.RandomBytes(16)
	// //encrypt the symmetric key with the public key
	// encrypted_sym_key, err_hybridencryption := userlib.PKEEnc(recipient_encryption_key, random_sym_key)

	// //store the encrypted symmetric key

	// if err_hybridencryption != nil{
	// 	err = errors.New("in encrypt_signature_invitationptr(): hybrid encryption for invitation_pointer failed")
	// 	return nil, err
	// }
	// //use the symmetric key to encrypt the actual data
	// iv := userlib.RandomBytes(16)
	// encrypted_invitationptr := userlib.SymEnc(random_sym_key, iv, invitationptr_byte)

	signature_invitationptr, err_signature := userlib.DSSign(caller_sign_key, encrypted_invitationptr)
	if err_signature != nil {
		err = errors.New("in encrypt_signature_invitationptr(): generate signature for invitatin_pointer failed")
		return nil, err
	}
	encrypted_invitationptr_final = append(encrypted_invitationptr, signature_invitationptr...)
	return encrypted_invitationptr_final, nil
}

// 用login public key来verify signature（integrity)
// 用输入password生成key，解码struct
func verify_decrypt_userstruct_signature(encrypted_userdata_final []byte, login_public_key userlib.DSVerifyKey, password string, username string) (userdata User, err error) {
	var fail_user User //如果获取失败，返回这个（好像默认是空的
	//把密码和加密了的userstruct从那一长串加密的东西中拖出来（分离出256位的signature）
	signature := encrypted_userdata_final[len(encrypted_userdata_final)-256:]
	encrypted_userstruct := encrypted_userdata_final[:len(encrypted_userdata_final)-256]

	//用login public key来verify signature（integrity)
	err_signature := userlib.DSVerify(login_public_key, encrypted_userstruct, signature)
	if err_signature != nil {
		err = errors.New("in verify_decrypt_userstruct_signature(): attacker modifies the content of user struct")
		return fail_user, err
	}

	////////////////////////////////用输入password生成key，解码struct
	//key是不是这样生成(symmetric_key)
	password_byte := []byte(password)
	salt_encrypt := []byte(username + "salt") //保证同一个密码的两个人不会same hash
	symmetric_key := userlib.Argon2Key(password_byte, salt_encrypt, 16)

	//如果ciphertext小于one block cipher，报错，不然就panic了
	if len(encrypted_userstruct) < userlib.AESBlockSizeBytes {
		err = errors.New("in verify_decrypt_userstruct_signature(): ciphertext input for SymDec shorter than one block cipher")
		return fail_user, err
	}
	decrypted_userstruct := userlib.SymDec(symmetric_key, encrypted_userstruct)

	err_marshal := json.Unmarshal(decrypted_userstruct, &userdata)
	if err_marshal != nil {
		err = errors.New("in verify_decrypt_userstruct_signature(): unmarshal deserilization failed")
		return fail_user, err
	}
	return userdata, nil
}

// 根据salt和password生成symmetric HMAC key和encryption key，对filestruct进行HMAC检查integrity，然后decrypt它
// ///////////////////////////////////////////////这是针对filestruct type1的
func verify_decrypt_filestruct_salt(encrypted_filestruct_final []byte, password string, username string, filename string) (decrypted_filestruct_byte []byte, err error) {
	//filestruct_fail := File_1{uuid.Nil, nil, nil, 0, uuid.Nil, nil, nil} //return value when failure
	//先取出filestruct
	correct_HMAC := encrypted_filestruct_final[len(encrypted_filestruct_final)-64:]
	encrypted_filestruct := encrypted_filestruct_final[:len(encrypted_filestruct_final)-64]
	//检查filestruct的integrity (deterministic的HMAC key)
	password_byte := []byte(password)
	salt_decrypt := []byte(username + filename + "cryptography") ///////////////////////////////要确定好salt是这样写
	salt_HMAC := []byte(username + filename + "HMAC")
	decrypt_symmetric_key := userlib.Argon2Key(password_byte, salt_decrypt, 16)
	HMAC_key_filestruct := userlib.Argon2Key(password_byte, salt_HMAC, 16)
	HMAC_filestruct, err_HMAC := userlib.HMACEval(HMAC_key_filestruct, encrypted_filestruct)
	if err_HMAC != nil {
		err = errors.New("in verify_decrypt_filestruct_salt(): the HMAC generation for file_struct checking failed")
		return nil, err
	}
	check := userlib.HMACEqual(HMAC_filestruct, correct_HMAC)
	if !check {
		err = errors.New("in verify_decrypt_filestruct_salt(): attacker modifies the file_struct")
		return nil, err
	}

	//解开filestruct

	//如果ciphertext小于one block cipher，报错，不然就panic了
	if len(encrypted_filestruct) < userlib.AESBlockSizeBytes {
		err = errors.New("in verify_decrypt_filestruct_salt(): ciphertext input for SymDec shorter than one block cipher")
		return nil, err
	}
	decrypted_filestruct_byte = userlib.SymDec(decrypt_symmetric_key, encrypted_filestruct)

	return decrypted_filestruct_byte, nil
}

// 根据提供的keys和uuid对filehead进行HMAC检查integrity，然后decrypt它
func verify_decrypt_filehead(encryption_key []byte, HMAC_key []byte, node_uuid uuid.UUID) (decrypted_filehead File_head, err error) {
	//获取file_head中能解开并检验file_content的symmetric key和HMAC key
	//获取加密了的filecontent+HMAC
	filehead_fail := File_head{uuid.Nil, nil, nil} //return value when failure
	encrypted_filenode_final, ok_datastore := userlib.DatastoreGet(node_uuid)
	if !ok_datastore { //找不到file content node
		err = errors.New("in verify_decrypt_filehead(): cannot find this file node in the datastore")
		return filehead_fail, err
	}
	correct_HMAC_filenode := encrypted_filenode_final[len(encrypted_filenode_final)-64:]
	encrypted_filenode := encrypted_filenode_final[:len(encrypted_filenode_final)-64]
	HMAC_filenode, err_HMAC := userlib.HMACEval(HMAC_key, encrypted_filenode)
	if err_HMAC != nil {
		err = errors.New("in verify_decrypt_filehead(): HMAC generation for filenode checking failed")
		//fmt.Println(err)
		return filehead_fail, err
	}
	check := userlib.HMACEqual(HMAC_filenode, correct_HMAC_filenode)
	if !check {
		err = errors.New("in verify_decrypt_filehead(): Attacker modifies the node in filenode, or user has no access to filehead")
		//fmt.Println(err)
		return filehead_fail, err
	}
	//解开file node
	var filehead File_head

	//如果ciphertext小于one block cipher，报错，不然就panic了
	if len(encrypted_filenode) < userlib.AESBlockSizeBytes {
		err = errors.New("in verify_decrypt_filehead(): ciphertext input for SymDec shorter than one block cipher")
		fmt.Println(err)
		return filehead_fail, err
	}
	decrypted_filenode_byte := userlib.SymDec(encryption_key, encrypted_filenode)

	err_marshal := json.Unmarshal(decrypted_filenode_byte, &filehead)
	if err_marshal != nil {
		err = errors.New("in verify_decrypt_filehead(): unmarshal deserilization for file_head failed")
		fmt.Println(err)
		return filehead_fail, err
	}

	return filehead, nil
}

// 根据提供的keys和uuid对descendentlist进行HMAC检查integrity，然后decrypt它
func verify_decrypt_descendentlist(encryption_key []byte, HMAC_key []byte, node_uuid uuid.UUID) (decrypted_descendentlist Descendent_list, err error) {
	//获取file_head中能解开并检验file_content的symmetric key和HMAC key
	//获取加密了的filecontent+HMAC
	descendent_list_fail := Descendent_list{nil} //return value when failure
	encrypted_descendentlist_final, ok_datastore := userlib.DatastoreGet(node_uuid)
	if !ok_datastore { //找不到descendent list
		err = errors.New("in verify_decrypt_descendentlist(): cannot find descendent list in the datastore")
		return descendent_list_fail, err
	}
	correct_HMAC_descendentlist := encrypted_descendentlist_final[len(encrypted_descendentlist_final)-64:]
	encrypted_descendentlist := encrypted_descendentlist_final[:len(encrypted_descendentlist_final)-64]
	HMAC_filenode, err_HMAC := userlib.HMACEval(HMAC_key, encrypted_descendentlist)
	if err_HMAC != nil {
		err = errors.New("in verify_decrypt_descendentlist(): HMAC generation for descendent list checking failed")
		return descendent_list_fail, err
	}
	check := userlib.HMACEqual(HMAC_filenode, correct_HMAC_descendentlist)
	if !check {
		err = errors.New("in verify_decrypt_descendentlist(): Attacker modifies the descendent list")
		return descendent_list_fail, err
	}
	//解开file node
	var descendent_list Descendent_list

	//如果ciphertext小于one block cipher，报错，不然就panic了
	if len(encrypted_descendentlist) < userlib.AESBlockSizeBytes {
		err = errors.New("in verify_decrypt_descendentlist(): ciphertext for SymDec shorter than one block cipher")
		return descendent_list_fail, err
	}
	decrypted_descendentlist_byte := userlib.SymDec(encryption_key, encrypted_descendentlist)

	err_marshal := json.Unmarshal(decrypted_descendentlist_byte, &descendent_list)
	if err_marshal != nil {
		err = errors.New("in verify_decrypt_descendentlist(): unmarshal deserilization for file_head failed")
		return descendent_list_fail, err
	}

	return descendent_list, nil
}

// 根据提供的keys和uuid对filecontent进行HMAC检查integrity，然后decrypt它
func verify_decrypt_filecontent(encryption_key []byte, HMAC_key []byte, node_uuid uuid.UUID) (decrypted_filecontent File_content, err error) {
	//获取file_head中能解开并检验file_content的symmetric key和HMAC key
	//获取加密了的filecontent+HMAC
	filecontent_fail := File_content{"", uuid.Nil, nil, nil}
	encrypted_filenode_final, ok_datastore := userlib.DatastoreGet(node_uuid)
	if !ok_datastore { //找不到file content node
		err = errors.New("in verify_decrypt_filecontent(): cannot find this file node in the datastore")
		return filecontent_fail, err
	}
	correct_HMAC_filenode := encrypted_filenode_final[len(encrypted_filenode_final)-64:]
	encrypted_filenode := encrypted_filenode_final[:len(encrypted_filenode_final)-64]
	HMAC_filenode, err_HMAC := userlib.HMACEval(HMAC_key, encrypted_filenode)
	if err_HMAC != nil {
		err = errors.New("in verify_decrypt_filecontent(): HMAC generation for filenode checking failed")
		return filecontent_fail, err
	}
	check := userlib.HMACEqual(HMAC_filenode, correct_HMAC_filenode)
	if !check {
		err = errors.New("in verify_decrypt_filecontent(): Attacker modifies the node in filenode")
		return filecontent_fail, err
	}
	//解开file node
	var filecontent File_content
	//如果ciphertext小于one block cipher，报错，不然就panic了
	if len(encrypted_filenode) < userlib.AESBlockSizeBytes {
		err = errors.New("in verify_decrypt_filecontent(): ciphertext input for SymDec shorter than one block cipher")
		return filecontent_fail, err
	}
	decrypted_filenode := userlib.SymDec(encryption_key, encrypted_filenode)

	err_marshal := json.Unmarshal(decrypted_filenode, &filecontent)
	if err_marshal != nil {
		err = errors.New("in verify_decrypt_filecontent(): unmarshal deserilization for file_head failed")
		return filecontent_fail, err
	}

	return filecontent, nil
}

// 根据提供的keys和uuid对filecontent进行HMAC检查integrity，然后decrypt它
func verify_decrypt_metadata(encryption_key []byte, HMAC_key []byte, node_uuid uuid.UUID) (decrypted_metadata Owner_share_metadata, err error) {
	//获取file_head中能解开并检验file_content的symmetric key和HMAC key
	//获取加密了的filecontent+HMAC
	metadata_fail := Owner_share_metadata{uuid.Nil, nil, nil}
	encrypted_filenode_final, ok_datastore := userlib.DatastoreGet(node_uuid)
	if !ok_datastore { //找不到metadata
		err = errors.New("in verify_decrypt_metadata(): cannot find this file node in the datastore")
		return metadata_fail, err
	}
	correct_HMAC_filenode := encrypted_filenode_final[len(encrypted_filenode_final)-64:]
	encrypted_filenode := encrypted_filenode_final[:len(encrypted_filenode_final)-64]
	HMAC_filenode, err_HMAC := userlib.HMACEval(HMAC_key, encrypted_filenode)
	if err_HMAC != nil {
		err = errors.New("in verify_decrypt_metadata(): HMAC generation for filenode checking failed")
		return metadata_fail, err
	}
	check := userlib.HMACEqual(HMAC_filenode, correct_HMAC_filenode)
	if !check {
		err = errors.New("in verify_decrypt_metadata(): Attacker modifies the node in filenode")
		return metadata_fail, err
	}
	//解开file node
	var metadata Owner_share_metadata

	//如果ciphertext小于one block cipher，报错，不然就panic了
	if len(encrypted_filenode) < userlib.AESBlockSizeBytes {
		err = errors.New("in verify_decrypt_metadata(): ciphertext input for SymDec shorter than one block cipher")
		return metadata_fail, err
	}
	decrypted_filenode := userlib.SymDec(encryption_key, encrypted_filenode)

	err_marshal := json.Unmarshal(decrypted_filenode, &metadata)
	if err_marshal != nil {
		err = errors.New("in verify_decrypt_metadata(): unmarshal deserilization for file_head failed")
		return metadata_fail, err
	}

	return metadata, nil
}

// 解密并验证invitation_ptr
// //用sender的sharing_verify key来verify integrity，用receiver的sharing_decryption key来解密invitation_ptr
func verify_decrypt_invitationptr_signature(encrypted_invitationptr_final []byte, senderUsername string, receiver_decryption_key userlib.PKEDecKey) (invitation_pointer Invitation_pointer, err error) {
	invitation_pointer_fail := Invitation_pointer{uuid.Nil, nil, nil}
	//把密码和加密了的userstruct从那一长串加密的东西中拖出来（分离出256位的signature）
	signature := encrypted_invitationptr_final[len(encrypted_invitationptr_final)-256:]
	encrypted_invitationptr := encrypted_invitationptr_final[:len(encrypted_invitationptr_final)-256]

	////用sender的sharing_verify key来verify integrity，用receiver的sharing_decryption key来解密invitation_ptr
	//把sender的sharing_verify key取出来
	keystore_key_sign := (senderUsername + "share_sign")
	sender_verify_key, ok_getkey := userlib.KeystoreGet(keystore_key_sign)
	if !ok_getkey {
		err = errors.New("in verify_decrypt_invitationptr_signature(): can't find the sender's sharing_verify key")
		return invitation_pointer_fail, err
	}
	err_signature := userlib.DSVerify(sender_verify_key, encrypted_invitationptr, signature)
	if err_signature != nil {
		err = errors.New("in verify_decrypt_invitationptr_signature(): attacker modifies the content of invitation pointer")
		return invitation_pointer_fail, err
	}
	decrypted_invitationptr, err_decrypt := userlib.PKEDec(receiver_decryption_key, encrypted_invitationptr)
	if err_decrypt != nil {
		err = errors.New("in verify_decrypt_invitationptr_signature(): decryption for invitation_ptr failed")
		return invitation_pointer_fail, err
	}

	err_unmarshal := json.Unmarshal(decrypted_invitationptr, &invitation_pointer)
	if err_unmarshal != nil {
		err = errors.New("in verify_decrypt_invitationptr_signature(): deserilization for invitation_ptr failed")
		return invitation_pointer_fail, err
	}
	return invitation_pointer, nil
}

// 检查file_head和所有的file_content的integrity
func verify_filehead_and_allcontent(filestruct_encryption_key []byte, filestruct_HMAC_key []byte, filehead_uuid uuid.UUID) (err error) {
	//file_head解密与检验//////////////////
	//获取解密了的filehead
	file_head, err := verify_decrypt_filehead(filestruct_encryption_key, filestruct_HMAC_key, filehead_uuid)
	if err != nil {
		return err
	}
	//获取file_head中能解开并检验file_content的symmetric key和HMAC key
	encryption_key := file_head.Encryption_key
	HMAC_key := file_head.HMAC_key
	node_uuid := file_head.Content_ptr

	//file_content解密与检验(要写一个while loop)//////////////////
	for {
		file_content, err := verify_decrypt_filecontent(encryption_key, HMAC_key, node_uuid)
		if err != nil {
			return err
		}
		encryption_key = file_content.Encryption_key
		HMAC_key = file_content.HMAC_key
		node_uuid = file_content.Previous_uuid
		if encryption_key == nil && HMAC_key == nil {
			break
		}
		if encryption_key == nil || HMAC_key == nil {
			err = errors.New("in verify_filehead_and_allcontent(): the encryption_Key or HMAC_key is nil, abonormal case")
			return err
		}

	}
	return nil
}

// 先解密byte（判断是type1还是type2），如果是type1，先验证descendent list，获取解filehead的key，如果是type2，就要多走一个node去metadata那里
func unmarshal_verify_getkey4filehead(decrypted_filestruct_byte []byte) (filestruct_encryption_key []byte, filestruct_HMAC_key []byte, filehead_uuid uuid.UUID, err error) {
	var file_struct_1 File_1
	var file_struct_2 File_2
	//先unmarshal成filestruct type 1试试看
	err_marshal := json.Unmarshal(decrypted_filestruct_byte, &file_struct_1)
	if err_marshal != nil {
		err = errors.New("in unmarshal_verify_getkey4filehead(): unmarshal deserilization for file_struct_1 failed")
		return nil, nil, uuid.Nil, err
	}

	file_type := file_struct_1.File_type

	//如果实际上是一个type2的file，重新unmarshal转一遍
	if file_type == 2 { //不是owner
		err_marshal := json.Unmarshal(decrypted_filestruct_byte, &file_struct_2)
		if err_marshal != nil {
			err = errors.New("in unmarshal_verify_getkey4filehead(): unmarshal deserilization for file_struct_2 failed")
			return nil, nil, uuid.Nil, err
		}

		//获取file_struct中能解开并检验file_head的symmetric key和HMAC key,以及file_head的uuid
		//不是owner,那就要走多一个node才能得到解开filehead的key
		filestruct2_encryption_key := file_struct_2.Encryption_key
		filestruct2_HMAC_key := file_struct_2.HMAC_key
		metadata_uuid := file_struct_2.Share_ptr
		//解开metadata
		decrypted_metadata, err := verify_decrypt_metadata(filestruct2_encryption_key, filestruct2_HMAC_key, metadata_uuid)
		if err != nil {
			return nil, nil, uuid.Nil, err
		}
		//获取metadata里面解file head的key
		filestruct_encryption_key = decrypted_metadata.Encryption_key
		filestruct_HMAC_key = decrypted_metadata.HMAC_key
		filehead_uuid = decrypted_metadata.Head_pointer

	} else { //是owner
		//owner还要重新检查一次descendent list
		//获取file_struct中能解开并检验descendent_list和的symmetric key和HMAC key,以及descendent_list的uuid
		descendentlist_encryption_key := file_struct_1.DescendentList_Encryption_key
		descendentlist_HMAC_key := file_struct_1.DescendentList_HMAC_key
		descendentlist_uuid := file_struct_1.DescendentList_ptr

		/////////////////////////////////////////descendent_list解密与检验//////////////////
		_, err = verify_decrypt_descendentlist(descendentlist_encryption_key, descendentlist_HMAC_key, descendentlist_uuid)
		if err != nil {
			return nil, nil, uuid.Nil, err
		}

		filestruct_encryption_key = file_struct_1.Encryption_key
		filestruct_HMAC_key = file_struct_1.HMAC_key
		filehead_uuid = file_struct_1.Head_ptr
	}

	return filestruct_encryption_key, filestruct_HMAC_key, filehead_uuid, nil

}

func InitUser(username string, password string) (userdataptr *User, err error) {
	var userdata User

	//username是否为空，error
	if username == "" {
		err = errors.New("empty username")
		return nil, err
	}

	//根据username生成datastore中uuid
	/////////////////////////////////////////////////////////一个有一个没有的情况还没写
	user_uuid, err_bytes := info_string2uuid(username)
	if err_bytes != nil {
		err = errors.New("generate uuid for username failed")
		return nil, err
	}

	//uuid如果存在，即username存在，error  //在dataStore和keystore两个都找
	_, ok_datastore := userlib.DatastoreGet(user_uuid)

	keystore_key_login := (username + "login")
	_, ok_keystore := userlib.KeystoreGet(keystore_key_login)

	if ok_keystore && !ok_datastore { //keystore中有，但datastore无
		err = errors.New("struct user is deleted by attackers")
		return nil, err
	}
	if ok_keystore && ok_datastore { //都找到了
		err = errors.New("the user already exist")
		return nil, err
	}
	if !ok_keystore && ok_datastore { //keystore中无，但datastore有
		err = errors.New("attacker creates a fake user struct")
		return nil, err
	}

	userdata.Username = username
	userdata.Password = password

	//存sharing_decrypt到user struct
	sharing_encrypt, sharing_decrypt, err_crypt := userlib.PKEKeyGen()
	if err_crypt != nil {
		err = errors.New("generate en/decryption key failed")
		return nil, err
	}
	userdata.Sharing_decrypt = sharing_decrypt

	//存sharing_sign到user struct
	sharing_sign, sharing_verify, err_signature := userlib.DSKeyGen()
	if err_signature != nil {
		err = errors.New("generate sign/verification key failed")
		return nil, err
	}
	userdata.Sharing_sign = sharing_sign

	//存（username+“share_encrypt”,encrypt public key)到keystore
	keystore_key_encrypt := (username + "share_encrypt")
	err_keystore_encrypt := userlib.KeystoreSet(keystore_key_encrypt, sharing_encrypt)
	if err_keystore_encrypt != nil {
		err = errors.New("store share encryption key in keystore failed")
		return nil, err
	}

	//存（username+“share_sign”,verify public key)到keystore
	keystore_key_sign := (username + "share_sign")
	err_keystore_sign := userlib.KeystoreSet(keystore_key_sign, sharing_verify)
	if err_keystore_sign != nil {
		err = errors.New("store share verify key in keystore failed")
		return nil, err
	}

	login_sign, login_verify, err_signature := userlib.DSKeyGen()
	if err_signature != nil {
		err = errors.New("generate sign/verification key failed")
		return nil, err
	}

	//存 （username+“login”,verify public key)到keystore
	err_keystore_login := userlib.KeystoreSet(keystore_key_login, login_verify) // keystore_key_login 即(username + "login")
	if err_keystore_login != nil {
		err = errors.New("store login verify key in keystore failed")
		return nil, err
	}

	//利用password和salt加密user struct
	///////////////////////////////marshal写法，参数是写userdata还是&userdata
	//随机生成一个iv  //hash之后的password作为key
	//通过修改salt （username->username++） 和同样的password来形成一个新的symmetric key，用来加密user struct
	encrypted_userdata_final, err := encrypt_userstruct_salt_signature(login_sign, password, username, userdata)
	if err != nil {
		return nil, err
	}
	//存(uuid, 加密后的userdata)进datastore
	userlib.DatastoreSet(user_uuid, encrypted_userdata_final)

	return &userdata, nil
}

func GetUser(username string, password string) (userdataptr *User, err error) {
	var userdata User
	userdataptr = &userdata
	//在keystore中取public key, 随便找一个key
	keystore_key_login := (username + "login")
	login_public_key, ok_keystore := userlib.KeystoreGet(keystore_key_login)

	//根据username生成datastore中uuid
	user_uuid, err_bytes := info_string2uuid(username)
	if err_bytes != nil {
		err = errors.New("generate uuid for username failed")
		return nil, err
	}

	//在datastore中取出加密了的userstruct+signature
	encrypted_userdata_final, ok_datastore := userlib.DatastoreGet(user_uuid)

	if ok_keystore && !ok_datastore { //keystore中有，但datastore无
		err = errors.New("struct user is deleted by attackers")
		return nil, err
	} else if !ok_keystore && !ok_datastore { //都找不到
		err = errors.New("the user does not exist")
		return nil, err
	} else if !ok_keystore && ok_datastore { //keystore中无，但datastore有
		err = errors.New("attacker creates a fake user struct")
		return nil, err
	} else { //都找得到，验证integrity和密码是否正确
		//把密码和加密了的userstruct从那一长串加密的东西中拖出来（分离出256位的signature）
		//用login public key来verify signature（integrity)
		////////////////////////////////用输入password生成key，解码userstruct
		userdata, err = verify_decrypt_userstruct_signature(encrypted_userdata_final, login_public_key, password, username)
		if err != nil {
			return nil, err
		}

		////////////////////////////////获取正确的密码，比对用户的密码输入是否正确
		decrypted_password := userdata.Password
		if decrypted_password != password {
			err = errors.New("User input wrong password")
			return nil, err
		}
		userdataptr := &userdata
		return userdataptr, nil
	}
}

func (userdata *User) StoreFile(filename string, content []byte) (err error) {
	var file_head File_head
	var file_content File_content
	var descendent_list Descendent_list
	username := userdata.Username
	password := userdata.Password
	filestruct_uuid, err_file := info_string2uuid(username + filename) //file_uuid由username和filename生成
	if err_file != nil {
		err = errors.New("generate uuid for filename failed")
		return err
	}

	//检查这个file是否存在
	encrypted_filestruct_final, ok_datastore := userlib.DatastoreGet(filestruct_uuid)

	if ok_datastore { //这个file在datastore中存在
		/////////////////////////////////////////file_struct的salt、password解密与检验//////////////////
		//获取解密了的filestruct_byte
		decrypted_filestruct_byte, err := verify_decrypt_filestruct_salt(encrypted_filestruct_final, password, username, filename)
		if err != nil {
			return err
		}

		//先unmarshal byte（判断是type1还是type2），如果是type1，先验证descendent list，获取解filehead的key，如果是type2，就要多走一个node去metadata那里
		filestruct_encryption_key, filestruct_HMAC_key, filehead_uuid, err := unmarshal_verify_getkey4filehead(decrypted_filestruct_byte)
		if err != nil {
			return err
		}

		//检查file_head和所有的file_content的integrity
		err = verify_filehead_and_allcontent(filestruct_encryption_key, filestruct_HMAC_key, filehead_uuid)
		if err != nil {
			return err
		}

		/////////////////////////////检验完毕，修改file_head(file_head已经取出来了）和file_content
		//新建一个file_content
		uuid_filecontent := uuid.New()
		file_content.Content = string(content) //////////////////////////string还是byte
		file_content.Previous_uuid = uuid.Nil
		file_content.Encryption_key = nil ////////////////////////null（因为它后面没东西了，不用解了）
		file_content.HMAC_key = nil       //////////////////////////null

		//更改file_head的内容
		file_head.Encryption_key = userlib.RandomBytes(16) /////////////////////////用新的key
		file_head.HMAC_key = userlib.RandomBytes(16)
		filehead_encryption_key := file_head.Encryption_key /////////////////////////用新的key来加密后面的file_content
		filehead_HMAC_key := file_head.HMAC_key
		file_head.Content_ptr = uuid_filecontent //改filehead指向的uuid值

		////////////////////////////////filestruct的key不变，拿取出来的这些key给新的file_content加密
		//file_head用新的key帮新的file_content加密
		encrypted_filecontent_final, err_filecontent := encrypt_HMAC_filenode(&file_content, filehead_encryption_key, filehead_HMAC_key)
		if err_filecontent != nil {
			err = errors.New("symmetric encryption and HMAC for new file_content failed")
			return err
		}

		//file_struct用原来的key帮新的file_head加密
		encrypted_filehead_final, err_filehead := encrypt_HMAC_filenode(&file_head, filestruct_encryption_key, filestruct_HMAC_key)
		if err_filehead != nil {
			err = errors.New("symmetric encryption and HMAC for file_head failed")
			return err
		}

		//把新的file_head(uuid不变）和file_content存进datastore
		userlib.DatastoreSet(filehead_uuid, encrypted_filehead_final)
		userlib.DatastoreSet(uuid_filecontent, encrypted_filecontent_final)
		return nil

	} else { //这个file在datastore中不存在
		//因为它第一次建，caller肯定是file_owner
		var file_struct File_1
		//file_content
		uuid_filecontent := uuid.New()
		file_content.Content = string(content) //////////////////////////////////////////string还是byte
		file_content.Previous_uuid = uuid.Nil
		file_content.Encryption_key = nil ////////////////////////null（因为它后面没东西了，不用解了）
		file_content.HMAC_key = nil       //////////////////////////null

		//file_head
		uuid_head := uuid.New()
		file_head.Content_ptr = uuid_filecontent
		file_head.Encryption_key = userlib.RandomBytes(16) //解file_content
		file_head.HMAC_key = userlib.RandomBytes(16)

		//这个file的descendent_list
		uuid_descendent_list := uuid.New()
		descendent_list.All_descendents = make(map[string]Descendent)

		//file_1 struct
		file_struct.Head_ptr = uuid_head
		file_struct.Encryption_key = userlib.RandomBytes(16)
		file_struct.HMAC_key = userlib.RandomBytes(16)
		file_struct.File_type = 1 //因为它第一次建，他肯定是file_owner
		file_struct.DescendentList_ptr = uuid_descendent_list
		file_struct.DescendentList_Encryption_key = userlib.RandomBytes(16) //用来解descendent_list的key
		file_struct.DescendentList_HMAC_key = userlib.RandomBytes(16)

		//////////////////////file_head帮file_content加密
		encrypted_filecontent_final, err_filecontent := encrypt_HMAC_filenode(&file_content, file_head.Encryption_key, file_head.HMAC_key)
		if err_filecontent != nil {
			err = errors.New("symmetric encryption and HMAC for file_content failed")
			return err
		}
		//////////////////////file_struct帮file_head加密
		encrypted_filehead_final, err_filehead := encrypt_HMAC_filenode(&file_head, file_struct.Encryption_key, file_struct.HMAC_key)
		if err_filehead != nil {
			err = errors.New("symmetric encryption and HMAC for file_head failed")
			return err
		}

		//////////////////////file_struct帮descendent_list加密
		encrypted_descendentlist_final, err_descendentlist := encrypt_HMAC_filenode(&descendent_list, file_struct.DescendentList_Encryption_key, file_struct.DescendentList_HMAC_key)
		if err_descendentlist != nil {
			err = errors.New("symmetric encryption and HMAC for descendent_list failed")
			return err
		}

		//存加密后的file_head,file_content和descendent_list进datastore
		userlib.DatastoreSet(uuid_filecontent, encrypted_filecontent_final)
		userlib.DatastoreSet(uuid_head, encrypted_filehead_final)
		userlib.DatastoreSet(uuid_descendent_list, encrypted_descendentlist_final)

		////////////////////用deterministic的symmetric key加密filestruct(因为是symmetric，所以不用存)
		//随机生成一个iv  //hash之后的password作为key
		//通过修改salt （username+filename+"__"） 和同样的password来形成一个新的symmetric key，用来加密file struct
		encrypted_filestruct_final, err := encrypt_filestruct_salt_HMAC(password, username, filename, file_struct)
		if err != nil {
			return err
		}

		//存(uuid, 加密后的filestruct)进datastore
		userlib.DatastoreSet(filestruct_uuid, encrypted_filestruct_final)

		return nil

	}

}

func (userdata *User) AppendToFile(filename string, content []byte) error {
	//append不用检验有没integrity问题
	var file_struct_1 File_1
	var file_struct_2 File_2
	var file_head File_head
	var new_node File_content
	var err error
	var filestruct_encryption_key, filestruct_HMAC_key []byte
	var filehead_uuid uuid.UUID
	username := userdata.Username
	password := userdata.Password
	filestruct_uuid, err_file := info_string2uuid(username + filename) //file_uuid由username和filename生成
	if err_file != nil {
		err = errors.New("generate uuid for filename failed")
		return err
	}

	//确定文件是否存在
	encrypted_filestruct_final, ok_datastore := userlib.DatastoreGet(filestruct_uuid)
	if !ok_datastore { //找不到文件，报错
		err := errors.New("file does not exist")
		return err
	}

	//获取解密了的filestruct_byte
	file_struct_byte, err := verify_decrypt_filestruct_salt(encrypted_filestruct_final, password, username, filename)
	if err != nil {
		return err
	}

	//////获取能解开filehead的key, 考虑到bandwidth问题，这里就不调用unmarshal_verify_getkey4filehead()了 //////////////////////////////////
	//////先解密byte（判断是type1还是type2），如果是type1，先验证descendent list，再直接获取解filehead的key
	//////如果是type2，就要多走一个node去metadata那里获取解filehead的key

	//先unmarshal成filestruct type 1试试看
	err_marshal := json.Unmarshal(file_struct_byte, &file_struct_1)
	if err_marshal != nil {
		err = errors.New("unmarshal deserilization into file_struct_1 failed")
		return err
	}

	file_type := file_struct_1.File_type

	//如果实际上是一个type2的file，重新unmarshal转一遍
	if file_type == 2 { //不是owner
		err_marshal := json.Unmarshal(file_struct_byte, &file_struct_2)
		if err_marshal != nil {
			err = errors.New("unmarshal deserilization into file_struct_2 failed")
			return err
		}

		//获取file_struct中能解开并检验file_head的symmetric key和HMAC key,以及file_head的uuid
		//不是owner,那就要走多一个node才能得到解开filehead的key
		filestruct2_encryption_key := file_struct_2.Encryption_key
		filestruct2_HMAC_key := file_struct_2.HMAC_key
		metadata_uuid := file_struct_2.Share_ptr
		//解开metadata
		decrypted_metadata, err := verify_decrypt_metadata(filestruct2_encryption_key, filestruct2_HMAC_key, metadata_uuid)
		if err != nil {
			return err
		}
		//获取metadata里面解file head的key
		filestruct_encryption_key = decrypted_metadata.Encryption_key
		filestruct_HMAC_key = decrypted_metadata.HMAC_key
		filehead_uuid = decrypted_metadata.Head_pointer

	} else { //是owner
		///////////////////考虑到append的bandwidth问题，owner就不重新检查descendent list了
		//直接获取file_struct中能解开并检验filehead的symmetric key和HMAC key,以及filehad的uuid
		filestruct_encryption_key = file_struct_1.Encryption_key
		filestruct_HMAC_key = file_struct_1.HMAC_key
		filehead_uuid = file_struct_1.Head_ptr
	}

	/////////////////////////////////////////file_head解密与检验//////////////////
	//获取解密了的filehead
	file_head, err = verify_decrypt_filehead(filestruct_encryption_key, filestruct_HMAC_key, filehead_uuid)
	if err != nil {
		return err
	}

	//获取file_head中能解开并检验node_1的symmetric key和HMAC key
	encryption_key := file_head.Encryption_key
	HMAC_key := file_head.HMAC_key
	node1_uuid := file_head.Content_ptr //取出file_head本来相连的第一个content（node1）的uuid

	//新建一个file_content(new_node)
	uuid_newnode := uuid.New()
	new_node.Content = string(content) //////////////////////////string还是byte
	new_node.Previous_uuid = node1_uuid
	new_node.Encryption_key = encryption_key ///////////////////////////newnode用filehead原来的那个key
	new_node.HMAC_key = HMAC_key

	//更改file_head的内容
	file_head.Encryption_key = userlib.RandomBytes(16) /////////////////////////filehead用新的key来加密new_node
	file_head.HMAC_key = userlib.RandomBytes(16)
	filehead_encryption_newkey := file_head.Encryption_key /////////////////////////用新的key来加密后面的file_content
	filehead_HMAC_newkey := file_head.HMAC_key
	file_head.Content_ptr = uuid_newnode //改filehead指向的uuid值

	//file_head用新的key帮new_node加密
	encrypted_newnode_final, err_filecontent := encrypt_HMAC_filenode(&new_node, filehead_encryption_newkey, filehead_HMAC_newkey)
	if err_filecontent != nil {
		err = errors.New("symmetric encryption and HMAC for new_node failed")
		return err
	}

	//file_struct用原来的key帮新的file_head加密
	encrypted_filehead_final, err_filehead := encrypt_HMAC_filenode(&file_head, filestruct_encryption_key, filestruct_HMAC_key)
	if err_filehead != nil {
		err = errors.New("symmetric encryption and HMAC for file_head failed")
		return err
	}

	//把新的file_head(uuid不变）和new_node存进datastore
	userlib.DatastoreSet(filehead_uuid, encrypted_filehead_final)
	userlib.DatastoreSet(uuid_newnode, encrypted_newnode_final)
	return nil

}

// ///////////////////////////////////////////the file may or may not be owned by caller
func (userdata *User) LoadFile(filename string) (content []byte, err error) {
	var file_head File_head
	var file_content File_content
	username := userdata.Username
	password := userdata.Password
	filestruct_uuid, err_file := info_string2uuid(username + filename) //file_uuid由username和filename生成
	if err_file != nil {
		err = errors.New("generate uuid for filename failed")
		return nil, err
	}
	encrypted_filestruct_final, ok_datastore := userlib.DatastoreGet(filestruct_uuid)
	if !ok_datastore { //找不到文件，报错
		err = errors.New("file does not exist")
		return nil, err
	}

	//获取解密了的filestruct
	file_struct_byte, err := verify_decrypt_filestruct_salt(encrypted_filestruct_final, password, username, filename)
	if err != nil {
		return nil, err
	}

	//先解密byte（判断是type1还是type2），
	//如果是type1，先验证descendent list，然后再直接获取解filehead的key
	//如果是type2，就要多走一个node去metadata那里获取解filehead的key
	filestruct_encryption_key, filestruct_HMAC_key, filehead_uuid, err := unmarshal_verify_getkey4filehead(file_struct_byte)
	if err != nil {
		return nil, err
	}

	//获取解密了的filehead
	file_head, err = verify_decrypt_filehead(filestruct_encryption_key, filestruct_HMAC_key, filehead_uuid)
	if err != nil {
		return nil, err
	}

	//获取file_head中能解开并检验file_content的symmetric key和HMAC key
	encryption_key := file_head.Encryption_key
	HMAC_key := file_head.HMAC_key
	node_uuid := file_head.Content_ptr

	var node_content []byte   //因为每一个节点的内容是string，它只能转成byte slice
	var content_list [][]byte //存所有节点读出来的内容，就要用一个二维列表了
	/////////////////////////////////////////file_content解密与检验(要写一个while loop)
	for {
		file_content, err = verify_decrypt_filecontent(encryption_key, HMAC_key, node_uuid)
		if err != nil {
			return nil, err
		}
		encryption_key = file_content.Encryption_key
		HMAC_key = file_content.HMAC_key
		node_uuid = file_content.Previous_uuid
		node_content = []byte(file_content.Content)
		content_list = append(content_list, node_content)
		if encryption_key == nil && HMAC_key == nil {
			break
		}
		if encryption_key == nil || HMAC_key == nil {
			err = errors.New("the encryption_Key or HMAC_key is nil, abonormal case")
			return nil, err
		}

	}
	reverse_content_list := reverse_flatten(content_list) //反转并flatten内容byte slice
	return reverse_content_list, nil

}

func (userdata *User) CreateInvitation(filename string, recipientUsername string) (invitationPtr uuid.UUID, err error) {
	var file_struct_1 File_1
	var file_struct_2 File_2
	var descendent_list Descendent_list
	var share_metadata Owner_share_metadata
	var invitation_pointer Invitation_pointer
	username := userdata.Username
	password := userdata.Password
	caller_sign_key := userdata.Sharing_sign
	filestruct_uuid, err_file := info_string2uuid(username + filename) //file_uuid由username和filename生成
	if err_file != nil {
		err = errors.New("generate uuid for filename in file existence check failed")
		return uuid.Nil, err
	}

	//检查caller所要share的文件是否在datastore中存在
	encrypted_filestruct_final, ok_datastore := userlib.DatastoreGet(filestruct_uuid)
	if !ok_datastore { //找不到文件，报错
		err = errors.New("file to be shared does not exist")
		return uuid.Nil, err
	}

	recipient_uuid, err_bytes := info_string2uuid(recipientUsername)
	if err_bytes != nil {
		err = errors.New("generate uuid for username in recipient name existence check failed")
		return uuid.Nil, err
	}

	//检查这个recipient是否存在
	_, ok_datastore = userlib.DatastoreGet(recipient_uuid)
	if !ok_datastore {
		err = errors.New("the recipient name does not exists")
		return uuid.Nil, err
	}

	////share之前检查这个file是否integrity
	//获取解密了的filestruct
	file_struct_byte, err := verify_decrypt_filestruct_salt(encrypted_filestruct_final, password, username, filename)
	if err != nil {
		return uuid.Nil, err
	}

	/////////////////////先unmarshal byte（判断是type1还是type2），如果是type1，先验证descendent list，获取解filehead的key
	/////////////////////如果是type2，就要多走一个node去metadata那里
	//先解密成filestruct type 1试试看
	err_marshal := json.Unmarshal(file_struct_byte, &file_struct_1)
	if err_marshal != nil {
		err = errors.New("unmarshal deserilization for file_struct_1 failed")
		return uuid.Nil, err
	}

	//file_struct(metadata)中能解开并检验file_head的symmetric key和HMAC key,以及file_head的uuid
	var filestruct_encryption_key []byte //用来解filehead的key
	var filestruct_HMAC_key []byte
	var filehead_uuid uuid.UUID

	file_type := file_struct_1.File_type

	//如果实际上是一个type2的file，重新转一遍
	if file_type == 2 { //不是owner
		err_marshal := json.Unmarshal(file_struct_byte, &file_struct_2)
		if err_marshal != nil {
			err = errors.New("unmarshal deserilization for file_struct_2 failed")
			return uuid.Nil, err
		}

		//获取file_struct中能解开并检验file_head的symmetric key和HMAC key,以及file_head的uuid
		//不是owner,那就要走多一个node才能得到解开filehead的key
		filestruct2_encryption_key := file_struct_2.Encryption_key
		filestruct2_HMAC_key := file_struct_2.HMAC_key
		metadata_uuid := file_struct_2.Share_ptr
		//解开metadata
		decrypted_metadata, err := verify_decrypt_metadata(filestruct2_encryption_key, filestruct2_HMAC_key, metadata_uuid)
		if err != nil {
			return uuid.Nil, err
		}
		//获取metadata里面解file head的key
		filestruct_encryption_key = decrypted_metadata.Encryption_key
		filestruct_HMAC_key = decrypted_metadata.HMAC_key
		filehead_uuid = decrypted_metadata.Head_pointer

		//检查file_head和所有的file_content的integrity
		err = verify_filehead_and_allcontent(filestruct_encryption_key, filestruct_HMAC_key, filehead_uuid)
		if err != nil {
			return uuid.Nil, err
		}

		////////以上检验完毕

		//不是owner就不用产生metadata，直接生成invitation pointer
		//新建invitation_ptr
		uuid_invitationptr := uuid.New()
		invitation_pointer.Meta = metadata_uuid
		invitation_pointer.Encrypt = filestruct2_encryption_key //就是filestruct拿来解metadata的key
		invitation_pointer.HMAC = filestruct2_HMAC_key

		////////////////加密invitation_ptr
		//用recipient的public share key去encrypt,用caller的private sign key去sign
		//recipient没接受，也依然用recipient的encrypt和sign
		//把recipient的public share key取出来
		keystore_key_encrypt := (recipientUsername + "share_encrypt")
		recipient_encryption_key, ok_getkey := userlib.KeystoreGet(keystore_key_encrypt)
		if !ok_getkey { //找不到sharing_encrypt_key
			err = errors.New("cannot find the public share encryption key for recipient in keystore")
			return uuid.Nil, err
		}
		encrypted_invitationptr_final, err := encrypt_signature_invitationptr(invitation_pointer, recipient_encryption_key, caller_sign_key, recipientUsername)
		if err != nil {
			return uuid.Nil, err
		}

		//把invitation_ptr存进datastore里面
		userlib.DatastoreSet(uuid_invitationptr, encrypted_invitationptr_final)

		return uuid_invitationptr, nil

	} else { //// 是owner
		//owner还要重新检查一次descendent list
		//获取file_struct中能解开并检验descendent_list和的symmetric key和HMAC key,以及descendent_list的uuid
		descendentlist_encryption_key := file_struct_1.DescendentList_Encryption_key
		descendentlist_HMAC_key := file_struct_1.DescendentList_HMAC_key
		descendentlist_uuid := file_struct_1.DescendentList_ptr

		/////////////////////////////////////////descendent_list解密与检验//////////////////
		descendent_list, err = verify_decrypt_descendentlist(descendentlist_encryption_key, descendentlist_HMAC_key, descendentlist_uuid)
		if err != nil {
			return uuid.Nil, err
		}

		filestruct_encryption_key = file_struct_1.Encryption_key
		filestruct_HMAC_key = file_struct_1.HMAC_key
		filehead_uuid = file_struct_1.Head_ptr

		//检查file_head和所有的file_content的integrity
		err = verify_filehead_and_allcontent(filestruct_encryption_key, filestruct_HMAC_key, filehead_uuid)
		if err != nil {
			return uuid.Nil, err
		}

		////////以上检验完毕///////////////////////////////////////////////////////

		//A每share一次，都会生成一个metadata的info和invitation_ptr
		//新建一个share_metada
		uuid_metadata := uuid.New() /////////////////////////////////////metadata的uuid是直接uuid.new()
		share_metadata.Head_pointer = filehead_uuid
		share_metadata.Encryption_key = filestruct_encryption_key /////////////metadata里面的key就是用来解密filehead的俩key
		share_metadata.HMAC_key = filestruct_HMAC_key

		//新建invitation_ptr
		uuid_invitationptr := uuid.New()
		invitation_pointer.Meta = uuid_metadata
		invitation_pointer.Encrypt = userlib.RandomBytes(16) /////////////这个key用来加密解码metadata，直接randombyte
		invitation_pointer.HMAC = userlib.RandomBytes(16)
		invitationptr_encryption_key := invitation_pointer.Encrypt
		invitationptr_HMAC_key := invitation_pointer.HMAC

		//invitation_ptr给metadata加密
		encrypted_metadata_final, err_metadata := encrypt_HMAC_filenode(&share_metadata, invitationptr_encryption_key, invitationptr_HMAC_key)
		if err_metadata != nil {
			err = errors.New("symmetric encryption and HMAC for metadata failed")
			return uuid.Nil, err
		}

		////////////////加密invitation_ptr
		//////用recipient的public share key去encrypt,用caller的private sign key去sign
		//////recipient没接受，也依然用recipient的encrypt和sign
		//把recipient的public share key取出来
		keystore_key_encrypt := (recipientUsername + "share_encrypt")
		recipient_encryption_key, ok_getkey := userlib.KeystoreGet(keystore_key_encrypt)
		if !ok_getkey { //找不到sharing_encrypt_key
			err = errors.New("cannot find the public share encryption key for recipient in keystore")
			return uuid.Nil, err
		}
		encrypted_invitationptr_final, err := encrypt_signature_invitationptr(invitation_pointer, recipient_encryption_key, caller_sign_key, recipientUsername)
		if err != nil {
			return uuid.Nil, err
		}

		//////////A的filestruct的descendent list，会加入一个新人
		var new_descendent Descendent
		new_descendent.Child_metadata_ptr = uuid_metadata
		new_descendent.Child_metadata_Encryption_key = invitationptr_encryption_key ///////////这是用来解密metadata的key,所以是和invitationptr同一个
		new_descendent.Child_metadata_HMAC_key = invitationptr_HMAC_key
		new_descendent.Invitationptr_uuid = uuid_invitationptr	//方便revoke的时候删库

		descendent_list.All_descendents[recipientUsername] = new_descendent

		//file_struct的key给descendent_list重新加密
		encrypted_descendentlist_final, err_descendentlist := encrypt_HMAC_filenode(&descendent_list, descendentlist_encryption_key, descendentlist_HMAC_key)
		if err_descendentlist != nil {
			err = errors.New("symmetric encryption and HMAC for descendent_list failed")
			return uuid.Nil, err
		}

		//把加密后的metadata, invitation ptr和descendent_list存到datastore里去
		userlib.DatastoreSet(uuid_metadata, encrypted_metadata_final)
		userlib.DatastoreSet(uuid_invitationptr, encrypted_invitationptr_final)
		userlib.DatastoreSet(descendentlist_uuid, encrypted_descendentlist_final)
		return uuid_invitationptr, nil
	}

}

func (userdata *User) AcceptInvitation(senderUsername string, invitationPtr uuid.UUID, filename string) error {
	var err error
	var invitation_pointer Invitation_pointer
	var file_struct_type2 File_2
	receiver_decryption_key := userdata.Sharing_decrypt
	username := userdata.Username
	password := userdata.Password
	filestruct_uuid, err_file := info_string2uuid(username + filename) //file_uuid由username和filename生成
	if err_file != nil {
		err = errors.New("generate uuid for filename in file existence check failed")
		return err
	}

	//检查接收者recipient所要起的这个filename是否已经在datastore中存在
	_, ok_datastore := userlib.DatastoreGet(filestruct_uuid)
	if ok_datastore { //文件已经存在
		err = errors.New("the chosen filename already exists")
		return err
	}

	sender_uuid, err_bytes := info_string2uuid(senderUsername)
	if err_bytes != nil {
		err = errors.New("generate uuid for username in sender name existence check failed")
		return err
	}

	//检查sender是否存在
	_, ok_datastore = userlib.DatastoreGet(sender_uuid)
	if !ok_datastore {
		err = errors.New("the sender does not exists")
		return err
	}

	//把加密了的invitation_ptr从keystore里面取出来
	encrypted_invitationptr_final, ok_datastore := userlib.DatastoreGet(invitationPtr)
	if !ok_datastore { //invitation_ptr不存在
		err = errors.New("the invitationptr content does not exists, or it has been invalid")
		return err
	}

	//解密并验证invitation_ptr
	//把密码和加密了的userstruct从那一长串加密的东西中拖出来（分离出256位的signature）
	invitation_pointer, err = verify_decrypt_invitationptr_signature(encrypted_invitationptr_final, senderUsername, receiver_decryption_key)
	if err != nil {
		return err
	}

	//把解密的invitatoin_ptr的有用东西取出来
	metadata_uuid := invitation_pointer.Meta
	encryption_key := invitation_pointer.Encrypt
	HMAC_key := invitation_pointer.HMAC

	//新建一个file_struct_type2,把invitation_ptr存的用来解密metadata的key跟metadata的uuid搬到file_struct_type2里面
	file_struct_type2.Share_ptr = metadata_uuid
	file_struct_type2.Encryption_key = encryption_key
	file_struct_type2.HMAC_key = HMAC_key
	file_struct_type2.File_type = 2
	file_struct_type2.Parent_name = senderUsername

	//加密file_struct_type2
	encrypted_filestruct2_final, err := encrypt_filestruct_salt_HMAC(password, username, filename, file_struct_type2)
	if err != nil {
		return err
	}

	//把file_struct_type2存进datastore
	userlib.DatastoreSet(filestruct_uuid, encrypted_filestruct2_final)

	return nil
}

func (userdata *User) RevokeAccess(filename string, recipientUsername string) error {
	var err error
	var file_head File_head
	var file_struct File_1 //肯定是file_1
	username := userdata.Username
	password := userdata.Password

	filestruct_uuid, err_file := info_string2uuid(username + filename) //file_uuid由username和filename生成
	if err_file != nil {
		err = errors.New("generate uuid for filename in file existence check failed")
		return err
	}

	////检查所要起的这个filename是否已经在datastore中存在
	encrypted_filestruct_final, ok_datastore := userlib.DatastoreGet(filestruct_uuid)
	if !ok_datastore { //所要revoke的文件不存在
		err = errors.New("the chosen filename does not exist")
		return err
	}

	recipient_uuid, err_bytes := info_string2uuid(recipientUsername)
	if err_bytes != nil {
		err = errors.New("generate uuid for username in recipient name existence check failed")
		return err
	}

	//检查这个recipient是否存在
	_, ok_datastore = userlib.DatastoreGet(recipient_uuid)
	if !ok_datastore {
		err = errors.New("the recipient name does not exists")
		return err
	}

	//先解密file_struct
	file_struct_byte, err := verify_decrypt_filestruct_salt(encrypted_filestruct_final, password, username, filename)
	if err != nil {
		return err
	}

	////先unmarshal byte，肯定是type1，因为caller一定是owner，先验证descendent list，获取解filehead的key
	//解密成filestruct type 1
	err_marshal := json.Unmarshal(file_struct_byte, &file_struct)
	if err_marshal != nil {
		err = errors.New("unmarshal deserilization for file_struct_1 failed")
		return err
	}

	//file_struct(metadata)中能解开并检验file_head的symmetric key和HMAC key,以及file_head的uuid
	var filestruct_encryption_key []byte //用来解filehead的key
	var filestruct_HMAC_key []byte
	var filehead_uuid uuid.UUID

	//owner还要重新检查一次descendent list
	//获取file_struct中能解开并检验descendent_list和的symmetric key和HMAC key,以及descendent_list的uuid
	descendentlist_encryption_key := file_struct.DescendentList_Encryption_key
	descendentlist_HMAC_key := file_struct.DescendentList_HMAC_key
	descendentlist_uuid := file_struct.DescendentList_ptr

	/////////////////////////////////////////descendent_list解密与检验//////////////////
	descendent_list, err := verify_decrypt_descendentlist(descendentlist_encryption_key, descendentlist_HMAC_key, descendentlist_uuid)
	if err != nil {
		return err
	}

	filestruct_encryption_key = file_struct.Encryption_key
	filestruct_HMAC_key = file_struct.HMAC_key
	filehead_uuid = file_struct.Head_ptr

	////检查file_head和所有的file_content的integrity
	err = verify_filehead_and_allcontent(filestruct_encryption_key, filestruct_HMAC_key, filehead_uuid)
	if err != nil {
		return err
	}

	/////////////////////////////////检验完毕

	//还是要file_head的
	file_head, err = verify_decrypt_filehead(filestruct_encryption_key, filestruct_HMAC_key, filehead_uuid)
	if err != nil {
		return err
	}

	//判断这个文件是否在share着(是否之前已经被revoke过一次了)
	_, in_map := descendent_list.All_descendents[recipientUsername]
	if !in_map {
		err = errors.New("the file is not currently shared with the recipient user, or it has been revoked")
		return err
	}

	////////////////////开始revoke access
	//把file_struct的用来解filehead的key改了
	file_struct.Encryption_key = userlib.RandomBytes(16)
	file_struct.HMAC_key = userlib.RandomBytes(16)

	//新的encryption_key和HMAC_KEY用来加密filehead
	new_filehead_encryption_key := file_struct.Encryption_key
	new_filehead_HMAC_key := file_struct.HMAC_key

	var metadata_uuid uuid.UUID

	////检查filename是否有跟recipient分享着
	revoked_flag := false	//是否已经revoke过了

	///////改没被revoke的descendent用来解密filehead的encryption key和HMAC key，以及他们对应的metadata (要遍历整个map)
	///////删掉要被revoke的那个key-value pair
	for descendent_name, descendent := range descendent_list.All_descendents {
		if descendent_name != recipientUsername {
			//一个循环内，对于descendent_name
			metadata_uuid = descendent.Child_metadata_ptr                       //他的metadata的地址
			metadata_Encryption_key := descendent.Child_metadata_Encryption_key ///////用来解密metadata的encryption key
			metadata_HMAC_key := descendent.Child_metadata_HMAC_key             ///////这些key不变？？？？

			//解密descendent的metadata
			decrypted_metadata, err := verify_decrypt_metadata(metadata_Encryption_key, metadata_HMAC_key, metadata_uuid)
			if err != nil {
				return err
			}

			//修改metadata里面相应的key
			decrypted_metadata.Encryption_key = new_filehead_encryption_key
			decrypted_metadata.HMAC_key = new_filehead_HMAC_key

			//重新加密这条metadata
			encrypted_metadata_new, err := encrypt_HMAC_filenode(&decrypted_metadata, metadata_Encryption_key, metadata_HMAC_key)
			if err != nil {
				return err
			}

			//重新储存新的metadata
			userlib.DatastoreSet(metadata_uuid, encrypted_metadata_new)

		} else { //从descendent_list中删除recipient_name
			revoked_flag = true
			target_descendent := descendent_list.All_descendents[descendent_name]
			invitationptr_uuid := target_descendent.Invitationptr_uuid
			userlib.DatastoreDelete(invitationptr_uuid)	//把这个descendent对应的invitation_ptr一并从datastore里面删掉
			delete(descendent_list.All_descendents, descendent_name)	//从descendent_list的这个descendent map里面删掉
		}
	}

	if !revoked_flag{
		err = errors.New("the file have been revoked from this user")
		return err
	}

	////重新、加密存储file_struct, filehead和descendent_list
	//用新的key重新加密file_head
	encrypted_filehead_new, err := encrypt_HMAC_filenode(&file_head, new_filehead_encryption_key, new_filehead_HMAC_key)
	if err != nil {
		return err
	}

	//重新加密删改过后的descendent_list
	encrypted_descendentlist_new, err := encrypt_HMAC_filenode(&descendent_list, descendentlist_encryption_key, descendentlist_HMAC_key)
	if err != nil {
		return err
	}

	//重新加密file_struct
	encrypted_filestruct_new, err := encrypt_filestruct_salt_HMAC(password, username, filename, file_struct)
	if err != nil {
		return err
	}

	//存储file_struct, filehead和descendent_list
	userlib.DatastoreSet(filestruct_uuid, encrypted_filestruct_new)
	userlib.DatastoreSet(filehead_uuid, encrypted_filehead_new)
	userlib.DatastoreSet(descendentlist_uuid, encrypted_descendentlist_new)

	return nil
}
