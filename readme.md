# EASY SHARE

EASY SHARE is a simple and efficient implementation in Go that allows you to share files from your execution directory easily. It leverages the Tor network to ensure that file transfers are anonymous, secure, and encrypted.

## Overview

EASY SHARE provides a straightforward way to share files without compromising your privacy. By using the Tor network, it masks your IP address and encrypts your data, making it safe to share sensitive information.

## How It Works

1. **Sharing Files**:
   - Execute the `share` binary.
   - Wait for at least 1 minute to establish a Tor connection.
   - Once connected, the console will display the links to the files available for sharing. These links are what you will send to the recipient.

2. **Downloading Files**:
   - The recipient should execute the `download` binary.
   - They must also wait for at least 1 minute to establish a Tor connection.
   - After the connection is established, the recipient should paste the received link into the console to begin the download.

## Instructions

### Sharing Files

#### Windows
1. Open Command Prompt and navigate to the directory containing the `share.exe`.
   - `cd path\to\your\directory`
   - `share.exe`
2. Wait for the console to display the file links. This may take a minute as the system establishes a connection to the Tor network.

#### Linux
1. Open a terminal and navigate to the directory containing the `share` binary.
   - `cd /path/to/your/directory`
   - `./share`
2. Wait for the console to display the file links. This may take a minute as the system establishes a connection to the Tor network.

### Receiving Files

#### Windows
1. Open Command Prompt and navigate to the directory containing the `download.exe`.
   - `cd path\to\your\directory`
   - `download.exe`
2. Wait for at least 1 minute to establish a Tor connection.
3. Paste the received link into the console to begin the download.

#### Linux
1. Open a terminal and navigate to the directory containing the `download` binary.
   - `cd /path/to/your/directory`
   - `./download`
2. Wait for at least 1 minute to establish a Tor connection.
3. Paste the received link into the console to begin the download.

## Important Notes

- Ensure that both the sender and receiver have the necessary permissions to execute the binaries.
- The Tor network may experience delays; please be patient during the connection process.
- For optimal performance, it is recommended to have a stable internet connection.

## Conclusion

EASY SHARE simplifies the process of file sharing while prioritizing user privacy and security. By utilizing the Tor network, you can confidently share files without the fear of being tracked or monitored.

For further inquiries or contributions, feel free to reach out!
 
