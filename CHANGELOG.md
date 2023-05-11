# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog][Keep a Changelog] and this project adheres to [Semantic Versioning][Semantic Versioning].

## [Unreleased]

## [v1.3.8] - 2021-11-27
### Added
- Added test for shared user being able to read appended file previously by the owner
- Added test for returning an error if a MITM attacker changes invitation
- Added test for writing and sharing a long files
- Added test for writing and sharing a file with empty content
- Added test for writing a file (and sharing) using empty string as a filename.
- Added test for creating and getting user with empty password

reasoning:
    For some reason, loading a file with empty content returned nil, instead of
    empty byte slice. Now, if length is 0, empty byte slice is manually returned.

## [v1.3.7] - 2021-11-25
### Added
- Added test for preventing access if inviter was revoked after invitation link was created

### Changed
- Changed RevokeAccess - now also deletes sharedList

## [v1.3.6] - 2021-11-25
### Added
- Added test for not being able to read file that does not exist
- Added test for not being able to append if access revoked
- Added test for not being able to create a new user with empty username

## [v1.3.5] - 2021-11-25
### Added
- Added test for user losing access to a file trying to write to this file

### Changed
- Changed GetFile to check kFile using key, then if failed, delete fileMetadata and return err
- Changed RevokeAccess - GenInvitationLink and delete are now put inside for loop

reasoning:
    Even though file key is deleted, metadata still remains, and this makes revoked
    user to think that they still have 'access' to the file.
    (They can't read it - it will return an error, but GetFile still returns exists)
    Owner of the file cannot remove this filemetadata, as the owner does not have access to
    kUser of sharer (or userdata).
    
    To solve this problem, when user querys file metadata, if it is found, the user will
    also try to get kfile using key in metadata. If it fails, this means that this user
    lost access to this file, so delete this metadata from datastore, then return err

## [v1.3.4] - 2021-11-24
### Added
- Added test for preventing access to files from attacks by stealing invitation links

## [v1.3.3] - 2021-11-24
### Removed
- Removed UpdateUserdata
- Removed UploadUserdata

## [v1.3.2] - 2021-11-24
### Changed
- Changed parameter of invLink to filePointer
- Now RevokeAccess removes filePointer stored at invLink
- Fixed a typo on test for not allowing regaining access if user's sharer's access got revoked

## [v1.3.1] - 2021-11-24
### Added
- Added tests for preventing access to files after revoking

## [v1.3.0] - 2021-11-23
### Changed
- Changed sharedList to now contain filePointer instead of usernames
- Changed tests in revoking access - removed typos

reasoning:
    GetFilePointer will return KFile, if this file is owned by the user.
    Since RevokeAccess will be called towards non-owner always,
    this filePointer should come from GetFile.
    (Using GetFilePointer is might seem ok to use, but it isn't due to reason below)
    
    GetFile (and GetFilePointer) is taking in userdata as parameter and
    is using kUser, filename and salt as key to datastore.
    If Alice shares Bob, and Bob shares to Charles,
    and if Alice revoke access Bob, Alice can successfully do so, but
    to remove access to Charles, Alice must have userdata of Bob.
    In RevokeAccess, Alice must remove Charles but Alice does not have access to
    Bob's userdata.
    
    So, instead of sharedList saving usernames, filepointer for those users are saved
    such that Alice has access to Charles filePointer and can delete kFile at
    that location.
    

## [v1.2.2] - 2021-11-23
### Added
- Added test for revoking access, as non-owner (who was shared by another user)

## [v1.2.1] - 2021-11-21
### Changed
- Changed indentation of some code

### Note
- Added to past commit to CHANGELOG

## [v1.2.0] - 2021-11-20
### Removed
- Removed Files from User struct

### Deprecated
- UpdateUserdata and UploadUserdata are deprecated, due to removing Files from User struct

### Added
- Added SetFile and GetFile

### Changed
- Changed all code referencing User.Files to SetFile and GetFile

reasoning:
    Since storing all file keys in User struct in form of slice slows down
    query as user has more files, Files are removed from User struct and
    have their own place in datastore, so that accessing one file will
    only require one small query

## [v1.1.3] - 2021-11-19
### Added
- Added test for revoking access not affecting other users

## [v1.1.2] - 2021-11-19
### Added
- Added test for revoking access with one user

## [v1.1.1] - 2021-11-19
### Added
- Added test for accepting and accessing file with different device

### Changed
- Used UploadUserdata in multiple methods to DRY the code
- Added UploadUserdata to CreateInvitation

reasoning:
    UpdateUserdata is good, but it only querys userdata to that in datastore.
    It never puts current userdata to datastore. (This was found in CreateInvitation)
    For example, if Alice logs into desktop and laptop, gives invitation link to Bob on
    her desktop, then RevokeAccess Bob with laptop?
    Files in AliceDesktop should be updated with Bob's username, but AliceLaptop
    would not have Bob's username, and cause error for revoking user who was not
    shared.

## [v1.1.1] - 2021-11-19
### Added
- Added UploadUserdata

## [v1.1.0] - 2021-11-19
### Added
- Added tests for multiple user sharing, which also passes

### Removed
- Removed StringSlice
- Removed MarshalSlice
- Removed UnmarshalSlice

### Changed
- AcceptInvitation now initialize sharedList to empty slice, but it does not pass yet

reasoning:
    Apparently, string slices can be marshalled, so StringSlice struct and its related methods
    are removed.
    CreateInvitation from Bob to Nilufar was causing error, because when Bob accepted invitation
    from Alice, sharedList for Bob's filepointer was not initialized to empty slice.
    This caused problem when Bob was trying to add Nilufar's username to sharedList,
    as sharedList did not exist in the first place.

## [v1.0.2] - 2021-11-19
### Added
- Added new test for append for both owner and shared user

### Changed
- Changed logic operating in AppendToFile

## [v1.0.1] - 2021-11-19
### Added
- Added UpdateUserdata

### Changed
- When userdata.Files is intended to be used, UpdateUserdata is called first, to make sure
    it is up to date

reasoning:
    CreateInvitation was failing, with error saying File is not found from user's account
    It was clearly passing our implemented test. Here were autograder's test:
    Alice logs into desktop
    Alice logs into laptop
    Alice creates file F on desktop
    Alice creates invitation link for Bob to share F on laptop

    On last step, it was having problem because when Alice creates F in desktop,
    User struct on laptop instance did not include F in user.Files

    To fix this, we had to query datastore every time user.Files is used

## [v1.0.1] - 2021-11-18
### Fixed
- Fixed GetUser (I don't remember how. I forgot)

## [v1.0.0] - 2021-11-18
### Added
- Most of features are added
- Due to working in local setting, CHANGELOG was not recorded

## [v0.2.0] - 2021-03-29
### Changed
- Updated [userlib][userlib] dependency to `v0.2.0`.

---

## [Released]

## [v0.1.0] - 2021-02-21
CHANGELOG did not exist in this release.

---

<!-- Links -->
[Keep a Changelog]: https://keepachangelog.com/
[Semantic Versioning]: https://semver.org/
[userlib]: https://github.com/cs161-staff/project2-userlib/blob/master/CHANGELOG.md

<!-- Versions -->
[Unreleased]: https://github.com/cs161-staff/project2-starter-code/compare/v0.2.0...HEAD
[Released]: https://github.com/cs161-staff/project2-starter-code/releases
[v0.2.0]: https://github.com/cs161-staff/project2-starter-code/compare/v0.1.0...v0.2.0
[v0.1.0]: https://github.com/cs161-staff/project2-starter-code/releases/v0.1.0
