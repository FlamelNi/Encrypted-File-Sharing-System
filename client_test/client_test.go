package client_test

// You MUST NOT change these default imports.  ANY additional imports it will
// break the autograder and everyone will be sad.

import (
	// Some imports use an underscore to prevent the compiler from complaining
	// about unused imports. Normally, you will want to avoid underscore imports
	// unless you know exactly what you are doing. You can read more about
	// underscore imports here: https://golangdocs.com/blank-identifier-in-golang
	_ "encoding/hex"
	_ "errors"
	_ "strconv"
	_ "strings"
	"testing"
	_ "fmt"

	// A "dot" import is used here so that the functions in the ginko and gomega
	// modules can be used without an identifier. For example, Describe() and
	// Expect() instead of ginko.Describe() and gomega.Expect(). You can read more
	// about dot imports here:
	// https://stackoverflow.com/questions/6478962/what-does-the-dot-or-period-in-a-go-import-statement-do
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	userlib "github.com/cs161-staff/project2-userlib"

	// The client implementation is intentionally defined in a different package.
	// This forces us to follow best practice and write tests that only rely on
	// client API that is exported from the client package, and avoid relying on
	// implementation details private to the client package.
	"github.com/cs161-staff/project2-starter-code/client"
)

func TestSetupAndExecution(t *testing.T) {
	// We are using 2 libraries to help us write readable and maintainable tests:
	//
	// (1) Ginkgo, a Behavior Driven Development (BDD) testing framework that
	//             makes it easy to write expressive specs that describe the
	//             behavior of your code in an organized manner; and
	//
	// (2) Gomega, an assertion/matcher library that allows us to write individual
	//             assertion statements in tests that read more like natural
	//             language. For example "Expect(ACTUAL).To(Equal(EXPECTED))".
	//
	// In the Ginko framework, a test case signals failure by calling Ginkgo’s
	// Fail(description string) function. However, we are using the Gomega library
	// to execute our assertion statements. When a Gomega assertion fails, Gomega
	// calls a GomegaFailHandler, which is a function that must be provided using
	// gomega.RegisterFailHandler(). Here, we pass Ginko's Fail() function to
	// Gomega so that Gomega can report failed assertions to the Ginko test
	// framework, which can take the appropriate action when a test fails.
	//
	// This is the sole connection point between Ginkgo and Gomega.
	RegisterFailHandler(Fail)

	RunSpecs(t, "Client Tests")
}

// ================================================
// Here are some optional global variables that can be used throughout the test
// suite to make the tests more readable and maintainable than defining these
// values in each test. You can add more variables here if you want and think
// they will help keep your code clean!
// ================================================
const someFilename = "file1.txt"
const someOtherFilename = "file2.txt"
const nonExistentFilename = "thisFileDoesNotExist.txt"

const aliceUsername = "Alice"
const alicePassword = "AlicePassword"
const bobUsername = "Bob"
const bobPassword = "BobPassword"
const nilufarUsername = "Nilufar"
const nilufarPassword = "NilufarPassword"
const olgaUsername = "Olga"
const olgaPassword = "OlgaPassword"
const marcoUsername = "Marco"
const marcoPassword = "MarcoPassword"

const nonExistentUsername = "NonExistentUser"

var alice *client.User
var bob *client.User
var nilufar *client.User
var olga *client.User
var marco *client.User

var someFileContent []byte
var someShortFileContent []byte
var someLongFileContent []byte

// ================================================
// The top level Describe() contains all tests in
// this test suite in nested Describe() blocks.
// ================================================

var _ = Describe("Client Tests", func() {
	BeforeEach(func() {
		// This top-level BeforeEach will be run before each test.
		//
		// Resets the state of Datastore and Keystore so that tests do not
		// interfere with each other.
		userlib.DatastoreClear()
		userlib.KeystoreClear()

		userlib.SymbolicDebug = false
		userlib.SymbolicVerbose = false
	})

	BeforeEach(func() {
		// This top-level BeforeEach will be run before each test.
		//
		// Byte slices cannot be constant, so this BeforeEach resets the content of
		// each global variable to a predefined value, which allows tests to rely on
		// the expected value of these variables.
		someShortFileContent = []byte("some short file content")
		someFileContent = someShortFileContent
		someLongFileContent = []byte("some LOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOONG file content")
	})

	Describe("Creating users", func() {
		It("should not error when creating a new user", func() {
			_, err := client.InitUser("Alice", "password")
			Expect(err).To(BeNil(), "Failed to initialized user Alice.")
		})
	
		It("should error if a username is already taken by another user", func() {
			_, err := client.InitUser("Alice", "password")
			_, err = client.InitUser("Alice", "password")
			Expect(err).ToNot(BeNil(), "Should error if a username is already taken")
		})
		
		It("should error when creating a new user with empty username", func() {
			_, err := client.InitUser("", "password")
			Expect(err).ToNot(BeNil(), "Successfully initialized user with empty username.")
		})

		It("should not error when creating a new user with empty password", func() {
			_, err := client.InitUser(aliceUsername, "")
			Expect(err).To(BeNil(), "Could not initialized a user with empty password.")
		})
	
		It("should error if a user does not exist with that username", func() {
			_, err := client.GetUser("Bob", "password")
			Expect(err).ToNot(BeNil(), "Should error because Bob has not been initialized")
		})
	
		// TODO: you probably want more test cases about creating users here
	})

	Describe("Getting users", func() {
		BeforeEach(func() {
			// This BeforeEach will run before each test in this Describe block.
			_, _ = client.InitUser("Alice", "some password")
		})
		
		It("should not error when getting a user", func() {
			_, err := client.GetUser("Alice", "some password")
			Expect(err).To(BeNil(), "Failed to log in user Alice.")
		})

		It("should not error when getting a user with empty password", func() {
			_, _ = client.InitUser(bobUsername, "")
			_, err := client.GetUser(bobUsername, "")
			Expect(err).To(BeNil(), "Failed to log in user Bob.")
		})
		
		It("should error when getting a user with wrong password", func() {
			_, err := client.GetUser("Alice", "wrong password")
			Expect(err).ToNot(BeNil(), "Should not log in user Alice.")
		})
		
		It("should error when getting a user with non-existent username", func() {
			_, err := client.GetUser("Nobdy", "some password")
			Expect(err).ToNot(BeNil(), "Should not log in to any user.")
		})
	})

	Describe("Single user storage", func() {
		var alice *client.User

		BeforeEach(func() {
			// This BeforeEach will run before each test in this Describe block.
			alice, _ = client.InitUser("Alice", "some password")
		})

		It("should upload content without erroring", func() {
			content := []byte("This is a test")
			err := alice.StoreFile("file1", content)
			Expect(err).To(BeNil(), "Failed to upload content to a file", err)
		})

		It("should download the expected content that was previously uploaded", func() {
			uploadedContent := []byte("This is a test")
			alice.StoreFile(someFilename, uploadedContent)
			downloadedContent, err := alice.LoadFile(someFilename)

			Expect(err).To(BeNil(), "Failed to download content from a file", err)
			
			Expect(downloadedContent).To(BeEquivalentTo(uploadedContent),
				"Downloaded content is not the same as uploaded content",
				downloadedContent,
				uploadedContent)
		})
		
		It("should not read file that does not exist", func() {
			_, err := alice.LoadFile(someFilename)
			Expect(err).ToNot(BeNil(), "Alice could read nonexistent file", err)
		})
		
		It("should error when trying to download a file that does not exist", func() {
			_, err := alice.LoadFile(nonExistentFilename)
			Expect(err).ToNot(BeNil(), "Was able to load a non-existent file without error.")
		})

		It("should not error when trying to append to a file", func() {
			content := []byte("This is a test")
			filename := "file1"
			err := alice.StoreFile(filename, content)
			Expect(err).To(BeNil(), "Failed to upload content to a file", err)

			err = alice.AppendToFile(filename, content)
			Expect(err).To(BeNil(), "Failed to append content to a file", err)

			downloadedContent, err := alice.LoadFile(filename)
			Expect(err).To(BeNil(), "Failed to download content from a file", err)
			
			Expect(downloadedContent).To(BeEquivalentTo(append(content, content...)),
				"Downloaded content is not the same as uploaded content")
		})

		It("should error when trying to append to a nonexistent file", func() {
			content := []byte("This is a test")
			filename := "file1"
			err := alice.StoreFile(filename, content)
			Expect(err).To(BeNil(), "Failed to upload content to a file", err)
		
			err = alice.AppendToFile("", content)
			Expect(err).ToNot(BeNil(), "Should error when appending to nonexistent file", err)
		
			downloadedContent, err := alice.LoadFile(filename)
			Expect(err).To(BeNil(), "Failed to download content from a file", err)
		
			Expect(downloadedContent).To(BeEquivalentTo(content),
				"Downloaded content is not the same as uploaded content")
		})
		
		It("should not error when trying to append to a file many times", func() {
			content := []byte("This is a test")
			var result []byte
			filename := "file1"

			result = append(result, content...)
			err := alice.StoreFile("file1", content)
			Expect(err).To(BeNil(), "Failed to upload content to a file", err)

			for i := 1; i < 1000; i++ {
				result = append(result, content...)
				err = alice.AppendToFile("file1", content)
				Expect(err).To(BeNil(), "Failed to append content to a file", err)
			}

			downloadedContent, err := alice.LoadFile(filename)
			Expect(err).To(BeNil(), "Failed to download content from a file", err)
			
			Expect(downloadedContent).To(BeEquivalentTo(result),
				"Downloaded content is not the same as uploaded content")
		})
	})

	Describe("Sharing files", func() {
		
		BeforeEach(func() {
			// Initialize each user to ensure the variable has the expected value for
			// the tests in this Describe() block.
			alice, _ = client.InitUser(aliceUsername, alicePassword)
			bob, _ = client.InitUser(bobUsername, bobPassword)
			nilufar, _ = client.InitUser(nilufarUsername, nilufarPassword)
			olga, _ = client.InitUser(olgaUsername, olgaPassword)
			marco, _ = client.InitUser(marcoUsername, marcoPassword)
		})

		It("should share a file without erroring", func() {
			alice.StoreFile(someFilename, someShortFileContent)
			shareFileInfoPtr, err := alice.CreateInvitation(someFilename, bobUsername)
			Expect(err).To(BeNil(), "Alice failed to share a file with Bob.")

			err = bob.AcceptInvitation(aliceUsername, shareFileInfoPtr, someOtherFilename)
			Expect(err).To(BeNil(), "Bob could not receive the file that Alice shared.")

			downloadedContent, err := bob.LoadFile(someOtherFilename)
			Expect(err).To(BeNil(), "Bob could not load the file that Alice shared.")
			Expect(downloadedContent).To(BeEquivalentTo(someShortFileContent),
				"The file contents that Bob downloaded was not the same as what Alice uploaded.")
		})

		It("should share a long file without erroring", func() {
			alice.StoreFile(someFilename, someLongFileContent)
			shareFileInfoPtr, err := alice.CreateInvitation(someFilename, bobUsername)
			Expect(err).To(BeNil(), "Alice failed to share a file with Bob.")

			err = bob.AcceptInvitation(aliceUsername, shareFileInfoPtr, someOtherFilename)
			Expect(err).To(BeNil(), "Bob could not receive the file that Alice shared.")

			downloadedContent, err := bob.LoadFile(someOtherFilename)
			Expect(err).To(BeNil(), "Bob could not load the file that Alice shared.")
			Expect(downloadedContent).To(BeEquivalentTo(someLongFileContent),
				"The file contents that Bob downloaded was not the same as what Alice uploaded.")
		})

		It("should share an empty file without erroring", func() {
			emptyFileContent := make([]byte, 0)
			alice.StoreFile(someFilename, emptyFileContent)
			shareFileInfoPtr, err := alice.CreateInvitation(someFilename, bobUsername)
			Expect(err).To(BeNil(), "Alice failed to share a file with Bob.")

			err = bob.AcceptInvitation(aliceUsername, shareFileInfoPtr, someOtherFilename)
			Expect(err).To(BeNil(), "Bob could not receive the file that Alice shared.")

			downloadedContent, err := bob.LoadFile(someOtherFilename)
			Expect(err).To(BeNil(), "Bob could not load the file that Alice shared.")
			Expect(downloadedContent).To(BeEquivalentTo(emptyFileContent),
				"The file contents that Bob downloaded was not the same as what Alice uploaded.")
		})

		It("should write, load and share a file with empty string as filename", func() {
			alice.StoreFile("", someFileContent)
			shareFileInfoPtr, err := alice.CreateInvitation("", bobUsername)
			Expect(err).To(BeNil(), "Alice failed to share a file with Bob.")
			
			err = bob.AcceptInvitation(aliceUsername, shareFileInfoPtr, "")
			Expect(err).To(BeNil(), "Bob could not receive the file that Alice shared.")
			
			downloadedContent, err := bob.LoadFile("")
			Expect(err).To(BeNil(), "Bob could not load the file that Alice shared.")
			Expect(downloadedContent).To(BeEquivalentTo(someFileContent),
				"The file contents that Bob downloaded was not the same as what Alice uploaded.")
			
			err = bob.AppendToFile("", []byte("hello"))
			Expect(err).To(BeNil(), "Bob should not get error when appending to file", err)
			
			err = alice.AppendToFile("", []byte("world"))
			Expect(err).To(BeNil(), "Alice should not get error when appending to file", err)
			
			downloadedContent, err = alice.LoadFile("")
			Expect(err).To(BeNil(), "Alice could not load the file that was appended.")
			Expect(downloadedContent).To(BeEquivalentTo(append(someFileContent, []byte("helloworld")...)),
				"The file contents that Alice downloaded was not the same as what Alice and Bob appended.")
		})
		
		It("should be able to accept and access file with different device", func() {
			alice.StoreFile(someFilename, someShortFileContent)
			shareFileInfoPtr, err := alice.CreateInvitation(someFilename, bobUsername)
			Expect(err).To(BeNil(), "Alice failed to share a file with Bob.")
			
			bobDesktop, err := client.GetUser(bobUsername, bobPassword)
			Expect(err).To(BeNil(), "Bob failed to log in in desktop.")
			
			bobLaptop, err := client.GetUser(bobUsername, bobPassword)
			Expect(err).To(BeNil(), "Bob failed to log in in laptop.")
			
			err = bobDesktop.AcceptInvitation(aliceUsername, shareFileInfoPtr, someOtherFilename)
			Expect(err).To(BeNil(), "Bob could not receive the file that Alice shared with desktop.")
			
			downloadedContent, err := bobLaptop.LoadFile(someOtherFilename)
			Expect(err).To(BeNil(), "Bob could not load the file that Alice shared with laptop.")
			Expect(downloadedContent).To(BeEquivalentTo(someShortFileContent),
				"The file contents that Bob downloaded was not the same as what Alice uploaded.")
		})
		
		It("should share a file with different device", func() {
			aliceLaptop, err := client.GetUser(aliceUsername, alicePassword)
			Expect(err).To(BeNil(), "Alice failed to log in in laptop.")

			aliceDesktop, err := client.GetUser(aliceUsername, alicePassword)
			Expect(err).To(BeNil(), "Alice failed to log in in desktop.")
			
			aliceLaptop.StoreFile(someFilename, someShortFileContent)
			shareFileInfoPtr, err := aliceDesktop.CreateInvitation(someFilename, bobUsername)
			Expect(err).To(BeNil(), "Alice failed to share a file with Bob.")

			err = bob.AcceptInvitation(aliceUsername, shareFileInfoPtr, someOtherFilename)
			Expect(err).To(BeNil(), "Bob could not receive the file that Alice shared.")

			downloadedContent, err := bob.LoadFile(someOtherFilename)
			Expect(err).To(BeNil(), "Bob could not load the file that Alice shared.")
			Expect(downloadedContent).To(BeEquivalentTo(someShortFileContent),
				"The file contents that Bob downloaded was not the same as what Alice uploaded.")
		})
		
		It("should let shared user to load and append", func() {
			alice.StoreFile(someFilename, someShortFileContent)
			shareFileInfoPtr, err := alice.CreateInvitation(someFilename, bobUsername)
			Expect(err).To(BeNil(), "Alice failed to share a file with Bob.")
		
			err = bob.AcceptInvitation(aliceUsername, shareFileInfoPtr, someOtherFilename)
			Expect(err).To(BeNil(), "Bob could not receive the file that Alice shared.")
		
			bob.StoreFile(someOtherFilename, someFileContent)
			Expect(err).To(BeNil(), "Bob could not write to the file that Alice shared.")
		
			downloadedContent, err := alice.LoadFile(someFilename)
			Expect(err).To(BeNil(), "Alice could not load the file that Bob shared.")
			Expect(downloadedContent).To(BeEquivalentTo(someFileContent),
				"The file contents that Alice downloaded was not the same as what Bob uploaded.")
		
			err = alice.AppendToFile(someFilename, []byte("hello"))
			Expect(err).To(BeNil(), "Should not error when appending to file", err)
			
			err = bob.AppendToFile(someOtherFilename, []byte("world"))
			Expect(err).To(BeNil(), "Should not error when appending to file", err)
			
			downloadedContent, err = alice.LoadFile(someFilename)
			Expect(err).To(BeNil(), "Alice could not load the file that was appended.")
			Expect(downloadedContent).To(BeEquivalentTo(append(someFileContent, []byte("helloworld")...)),
				"The file contents that Alice downloaded was not the same as what Alice and Bob appended.")
			
			downloadedContent, err = bob.LoadFile(someOtherFilename)
			Expect(err).To(BeNil(), "Bob could not load the file that was appended.")
			Expect(downloadedContent).To(BeEquivalentTo(append(someFileContent, []byte("helloworld")...)),
			"The file contents that Bob downloaded was not the same as what Alice and Bob appended.")
		})
		
		It("should let shared user to read an appended file", func() {
			content := []byte("This is a test")
			err := alice.StoreFile(someFilename, content)
			Expect(err).To(BeNil(), "Failed to upload content to a file", err)
			
			err = alice.AppendToFile(someFilename, content)
			Expect(err).To(BeNil(), "Failed to append content to a file", err)
			
			downloadedContent, err := alice.LoadFile(someFilename)
			Expect(err).To(BeNil(), "Failed to download content from a file", err)
			Expect(downloadedContent).To(BeEquivalentTo(append(content, content...)),
			"Downloaded content is not the same as uploaded content")
			
			shareFileInfoPtr, err := alice.CreateInvitation(someFilename, bobUsername)
			Expect(err).To(BeNil(), "Alice failed to share a file with Bob.")
			
			err = bob.AcceptInvitation(aliceUsername, shareFileInfoPtr, someOtherFilename)
			Expect(err).To(BeNil(), "Bob could not receive the file that Alice shared.")
			
			downloadedContent, err = bob.LoadFile(someOtherFilename)
			Expect(err).To(BeNil(), "Bob failed to download content from a file", err)
			Expect(downloadedContent).To(BeEquivalentTo(append(content, content...)),
			"Bob downloaded content different from uploaded content")
		})
		
	})

	Describe("Sharing files with multiple people", func() {
		BeforeEach(func() {
			// Initialize each user to ensure the variable has the expected value for
			// the tests in this Describe() block.
			alice, _ = client.InitUser(aliceUsername, alicePassword)
			bob, _ = client.InitUser(bobUsername, bobPassword)
			nilufar, _ = client.InitUser(nilufarUsername, nilufarPassword)
			olga, _ = client.InitUser(olgaUsername, olgaPassword)
			marco, _ = client.InitUser(marcoUsername, marcoPassword)
		})
		
		
		It("should share a file with multiple people without erroring", func() {
			alice.StoreFile(someFilename, someShortFileContent)
			shareFileInfoPtr, err := alice.CreateInvitation(someFilename, bobUsername)
			Expect(err).To(BeNil(), "Alice failed to share a file with Bob.")

			err = bob.AcceptInvitation(aliceUsername, shareFileInfoPtr, someOtherFilename)
			Expect(err).To(BeNil(), "Bob could not receive the file that Alice shared.")

			shareFileInfoPtr, err = alice.CreateInvitation(someFilename, nilufarUsername)
			Expect(err).To(BeNil(), "Alice failed to share a file with Nilufar.")

			err = nilufar.AcceptInvitation(aliceUsername, shareFileInfoPtr, someOtherFilename)
			Expect(err).To(BeNil(), "Nilufar could not receive the file that Alice shared.")
			
			shareFileInfoPtr, err = bob.CreateInvitation(someOtherFilename, olgaUsername)
			Expect(err).To(BeNil(), "Bob failed to share a file with Olga.")
			
			err = olga.AcceptInvitation(bobUsername, shareFileInfoPtr, someOtherFilename)
			Expect(err).To(BeNil(), "Olga could not receive the file that Bob shared.")
			
			downloadedContent, err := bob.LoadFile(someOtherFilename)
			Expect(err).To(BeNil(), "Bob could not load the file that Alice shared.")
			Expect(downloadedContent).To(BeEquivalentTo(someShortFileContent),
				"The file contents that Bob downloaded was not the same as what Alice uploaded.")
			downloadedContent, err = nilufar.LoadFile(someOtherFilename)
			Expect(err).To(BeNil(), "Nilufar could not load the file that Alice shared.")
			Expect(downloadedContent).To(BeEquivalentTo(someShortFileContent),
				"The file contents that Nilufar downloaded was not the same as what Alice uploaded.")
			downloadedContent, err = olga.LoadFile(someOtherFilename)
			Expect(err).To(BeNil(), "Olga could not load the file that Bob shared.")
			Expect(downloadedContent).To(BeEquivalentTo(someShortFileContent),
				"The file contents that Olga downloaded was not the same as what Alice uploaded.")
		})
	})
	
	Describe("Revoking access to files", func() {
		BeforeEach(func() {
			// Initialize each user to ensure the variable has the expected value for
			// the tests in this Describe() block.
			alice, _ = client.InitUser(aliceUsername, alicePassword)
			bob, _ = client.InitUser(bobUsername, bobPassword)
			nilufar, _ = client.InitUser(nilufarUsername, nilufarPassword)
			olga, _ = client.InitUser(olgaUsername, olgaPassword)
			marco, _ = client.InitUser(marcoUsername, marcoPassword)
			
			// Alice -> Bob -> Nilufar -> Olga
			// 		-> Marco
			alice.StoreFile(someFilename, someShortFileContent)
			
			shareFileInfoPtr, _ := alice.CreateInvitation(someFilename, bobUsername)
			_ = bob.AcceptInvitation(aliceUsername, shareFileInfoPtr, someOtherFilename)
			
			shareFileInfoPtr, _ = alice.CreateInvitation(someFilename, marcoUsername)
			_ = marco.AcceptInvitation(aliceUsername, shareFileInfoPtr, someOtherFilename)
			
			shareFileInfoPtr, _ = bob.CreateInvitation(someOtherFilename, nilufarUsername)
			_ = nilufar.AcceptInvitation(bobUsername, shareFileInfoPtr, someOtherFilename)
			
			shareFileInfoPtr, _ = nilufar.CreateInvitation(someFilename, olgaUsername)
			_ = olga.AcceptInvitation(nilufarUsername, shareFileInfoPtr, someOtherFilename)
		})
		
		It("should revoke access to a file without erroring", func() {
			err := alice.RevokeAccess(someFilename, marcoUsername)
			Expect(err).To(BeNil(), "Alice failed to revoke access to a file with Marco.")
			
			_, err = marco.LoadFile(someOtherFilename)
			Expect(err).ToNot(BeNil(), "Marco could load the file that Alice shared.")
		})
		
		It("should not be able to append if access revoked", func() {
			err := alice.RevokeAccess(someFilename, marcoUsername)
			Expect(err).To(BeNil(), "Alice failed to revoke access to a file with Marco.")
			
			err = marco.AppendToFile(someOtherFilename, []byte("some content"))
			Expect(err).ToNot(BeNil(), "Marco could append the file that Alice shared.")
		})
		
		It("should revoke access to a file without interrupting other sharer", func() {
			err := alice.RevokeAccess(someFilename, marcoUsername)
			Expect(err).To(BeNil(), "Alice failed to revoke access to a file with Marco.")
			
			downloadedContent, err := bob.LoadFile(someOtherFilename)
			Expect(err).To(BeNil(), "Bob could not load the file that Alice shared.")
			Expect(downloadedContent).To(BeEquivalentTo(someShortFileContent),
				"The file contents that Bob downloaded was not the same as what Alice uploaded.")
		})
		
		It("should revoke file access to a user, and all other users that were shared by this user (Owner)", func() {
			err := alice.RevokeAccess(someFilename, bobUsername)
			Expect(err).To(BeNil(), "Alice failed to revoke access to a file with Bob.")
			
			_, err = bob.LoadFile(someOtherFilename)
			Expect(err).ToNot(BeNil(), "Bob could load the file that Alice shared.")
			
			_, err = nilufar.LoadFile(someOtherFilename)
			Expect(err).ToNot(BeNil(), "Nilufar could load the file that Bob shared.")
			
			_, err = olga.LoadFile(someOtherFilename)
			Expect(err).ToNot(BeNil(), "Olga could load the file that Nilufar shared.")
		})
		
		It("should revoke file access to a user, and all other users that were shared by this user (Non-owner)", func() {
			err := bob.RevokeAccess(someOtherFilename, nilufarUsername)
			Expect(err).To(BeNil(), "Bob failed to revoke access to a file with Nilufar.")
			
			_, err = nilufar.LoadFile(someOtherFilename)
			Expect(err).ToNot(BeNil(), "Nilufar could load the file that Bob shared.")
			
			_, err = olga.LoadFile(someOtherFilename)
			Expect(err).ToNot(BeNil(), "Olga could load the file that Nilufar shared.")
		})
		
		
	})
	
	Describe("Preventing access to files after revoking", func() {
		BeforeEach(func() {
			// Initialize each user to ensure the variable has the expected value for
			// the tests in this Describe() block.
			alice, _ = client.InitUser(aliceUsername, alicePassword)
			bob, _ = client.InitUser(bobUsername, bobPassword)
			nilufar, _ = client.InitUser(nilufarUsername, nilufarPassword)
			olga, _ = client.InitUser(olgaUsername, olgaPassword)
			marco, _ = client.InitUser(marcoUsername, marcoPassword)
			
			alice.StoreFile(someFilename, someShortFileContent)
		})
		
		It("should not allow regaining access with previous invitation link after revoke", func() {
			shareFileInfoPtr, err := alice.CreateInvitation(someFilename, bobUsername)
			Expect(err).To(BeNil(), "Alice failed to share a file with Bob.")
			
			err = bob.AcceptInvitation(aliceUsername, shareFileInfoPtr, someOtherFilename)
			Expect(err).To(BeNil(), "Bob could not receive the file that Alice shared.")
			
			err = alice.RevokeAccess(someFilename, bobUsername)
			Expect(err).To(BeNil(), "Alice failed to revoke access to a file with Bob.")
			
			err = bob.AcceptInvitation(aliceUsername, shareFileInfoPtr, someOtherFilename)
			Expect(err).ToNot(BeNil(), "Bob could accept invitation again")
			
			_, err = bob.LoadFile(someOtherFilename)
			Expect(err).ToNot(BeNil(), "Bob could load the file that Alice shared.")
		})
		
		It("should not allow gaining access, if invitation was given, but revoked right after", func() {
			shareFileInfoPtr, err := alice.CreateInvitation(someFilename, bobUsername)
			Expect(err).To(BeNil(), "Alice could not share a file with Bob")
			
			err = alice.RevokeAccess(someFilename, bobUsername)
			Expect(err).To(BeNil(), "Alice failed to revoke access to a file with Bob.")
			
			err = bob.AcceptInvitation(aliceUsername, shareFileInfoPtr, someOtherFilename)
			Expect(err).ToNot(BeNil(), "Bob could accept invitation again")
			
			_, err = bob.LoadFile(someOtherFilename)
			Expect(err).ToNot(BeNil(), "Bob could load the file that Alice shared.")
		})
		
		It("should not allow regaining access if user's sharer's access got revoked", func() {
			shareFileInfoPtr, err := alice.CreateInvitation(someFilename, bobUsername)
			Expect(err).To(BeNil(), "Alice could not share a file with Bob")
			
			err = bob.AcceptInvitation(aliceUsername, shareFileInfoPtr, someOtherFilename)
			Expect(err).To(BeNil(), "Bob could not accept invitation given by Alice")
			
			shareFileInfoPtr, err = bob.CreateInvitation(someOtherFilename, nilufarUsername)
			Expect(err).To(BeNil(), "Bob could not share a file with Nilufar")
			
			err = nilufar.AcceptInvitation(bobUsername, shareFileInfoPtr, someOtherFilename)
			Expect(err).To(BeNil(), "Nilufar could not accept invitation given by Bob")
			
			err = alice.RevokeAccess(someFilename, bobUsername)
			Expect(err).To(BeNil(), "Alice failed to revoke access to a file with Bob.")
			
			_, err = nilufar.LoadFile(someOtherFilename)
			Expect(err).ToNot(BeNil(), "Nilufar could load the file that Bob shared.")
			
			err = nilufar.AcceptInvitation(bobUsername, shareFileInfoPtr, someOtherFilename)
			Expect(err).ToNot(BeNil(), "Nilufar could accept invitation again (with Bob username)")
			
			err = nilufar.AcceptInvitation(aliceUsername, shareFileInfoPtr, someOtherFilename)
			Expect(err).ToNot(BeNil(), "Nilufar could accept invitation again (with Alice username)")
		})
		
		It("should write to seperate file after access is revoked", func() {
			shareFileInfoPtr, err := alice.CreateInvitation(someFilename, bobUsername)
			Expect(err).To(BeNil(), "Alice failed to share a file with Bob.")
			
			err = bob.AcceptInvitation(aliceUsername, shareFileInfoPtr, someOtherFilename)
			Expect(err).To(BeNil(), "Bob could not receive the file that Alice shared.")
			
			err = alice.RevokeAccess(someFilename, bobUsername)
			Expect(err).To(BeNil(), "Alice failed to revoke access to a file with Bob.")
			
			err = bob.StoreFile(someOtherFilename, someLongFileContent)
			Expect(err).To(BeNil(), "Bob could not create and write to new file.")
			
			downloadedContent, err := bob.LoadFile(someOtherFilename)
			Expect(err).To(BeNil(), "Bob could not load the file that Bob wrote.")
			Expect(downloadedContent).To(BeEquivalentTo(someLongFileContent),
				"Downloaded content was different from what Bob uploaded",
				downloadedContent,
				someLongFileContent)
			
			downloadedContent, err = alice.LoadFile(someFilename)
			Expect(err).To(BeNil(), "Alice could not load the file that Alice wrote.")
			Expect(downloadedContent).To(BeEquivalentTo(someShortFileContent),
			"Downloaded content was different from what Alice uploaded",
				downloadedContent,
				someShortFileContent)
			
		})
		
		It("should prevent access if inviter was revoked after invitation link was created", func() {
			shareFileInfoPtr, err := alice.CreateInvitation(someFilename, bobUsername)
			Expect(err).To(BeNil(), "Alice failed to share a file with Bob.")
			
			err = bob.AcceptInvitation(aliceUsername, shareFileInfoPtr, someOtherFilename)
			Expect(err).To(BeNil(), "Bob could not receive the file that Alice shared.")
			
			shareFileInfoPtr, err = bob.CreateInvitation(someOtherFilename, nilufarUsername)
			Expect(err).To(BeNil(), "Bob failed to share a file with Nilufar.")
			
			err = alice.RevokeAccess(someFilename, bobUsername)
			Expect(err).To(BeNil(), "Alice failed to revoke access to a file with Bob.")
			
			err = nilufar.AcceptInvitation(bobUsername, shareFileInfoPtr, someOtherFilename)
			Expect(err).ToNot(BeNil(), "Nilufar could receive the file that Bob shared.")
			
			_, err = nilufar.LoadFile(someOtherFilename)
			Expect(err).ToNot(BeNil(), "Nilufar could load the file that Bob shared.")
		})
	})
	
	Describe("Preventing access to files from attacks", func() {
		BeforeEach(func() {
			// Initialize each user to ensure the variable has the expected value for
			// the tests in this Describe() block.
			alice, _ = client.InitUser(aliceUsername, alicePassword)
			bob, _ = client.InitUser(bobUsername, bobPassword)
			nilufar, _ = client.InitUser(nilufarUsername, nilufarPassword)
			olga, _ = client.InitUser(olgaUsername, olgaPassword)
			marco, _ = client.InitUser(marcoUsername, marcoPassword)
			
			alice.StoreFile(someFilename, someShortFileContent)
		})
		
		It("should not allow gaining access, if invitation link was stolen before used", func() {
			shareFileInfoPtr, err := alice.CreateInvitation(someFilename, bobUsername)
			Expect(err).To(BeNil(), "Alice could not share a file with Bob")
			
			err = marco.AcceptInvitation(aliceUsername, shareFileInfoPtr, someOtherFilename)
			Expect(err).ToNot(BeNil(), "Marco could accept invitation again")
			
			_, err = marco.LoadFile(someOtherFilename)
			Expect(err).ToNot(BeNil(), "Marco could load the file that Alice shared.")
		})
		
		It("should error if a MITM attacker changes invitation", func() {
			shareFileInfoPtr, err := alice.CreateInvitation(someFilename, bobUsername)
			Expect(err).To(BeNil(), "Alice could not share a file with Bob")
			
			marco.StoreFile(someFilename, someShortFileContent)
			otherShareFileInfoPtr, err := marco.CreateInvitation(someFilename, bobUsername)
			Expect(err).To(BeNil(), "Marco could not share a file with Bob")
			
			jsonData, ok := userlib.DatastoreGet(otherShareFileInfoPtr)
			Expect(ok).To(BeTrue(), "Marco invitation was not found")
			
			userlib.DatastoreSet(shareFileInfoPtr, jsonData)
			
			err = bob.AcceptInvitation(aliceUsername, shareFileInfoPtr, someOtherFilename)
			Expect(err).ToNot(BeNil(), "Bob accepted a botched invitation")
		})
	})
	

	// TODO: you probably want more Describe() blocks to contain tests related to
	//       logical test groupings other than the ones suggested above

})
