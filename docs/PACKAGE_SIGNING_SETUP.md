# Multi-Platform Package Signing Configuration

This document describes the GitHub Actions secrets required for building and signing multi-platform packages.

## Required GitHub Secrets

### Linux Package Signing

#### RPM Packages
- `RPM_SIGNING_KEY`: Base64-encoded GPG private key for signing RPM packages
  ```bash
  gpg --export-secret-keys --armor KEY_ID | base64 -w 0
  ```

#### DEB Packages
- `DEB_SIGNING_KEY`: Base64-encoded GPG private key for signing DEB packages
  ```bash
  gpg --export-secret-keys --armor KEY_ID | base64 -w 0
  ```

#### AppImage Packages
- `APPIMAGE_SIGNING_KEY`: Base64-encoded GPG private key for signing AppImage packages
  ```bash
  gpg --export-secret-keys --armor KEY_ID | base64 -w 0
  ```

### Windows Package Signing

#### Code Signing Certificate
- `WINDOWS_CERT_BASE64`: Base64-encoded PKCS#12 certificate file (.p12)
  ```bash
  cat certificate.p12 | base64 -w 0
  ```
- `WINDOWS_CERT_PASSWORD`: Password for the PKCS#12 certificate

### macOS Package Signing

#### Developer Certificate
- `MACOS_CERT_BASE64`: Base64-encoded PKCS#12 certificate file (.p12)
  ```bash
  cat certificate.p12 | base64 -w 0
  ```
- `MACOS_CERT_PASSWORD`: Password for the PKCS#12 certificate
- `MACOS_KEYCHAIN_PASSWORD`: Password for the temporary keychain
- `MACOS_DEVELOPER_ID`: Developer ID for code signing (e.g., "Developer ID Application: Your Name (TEAM_ID)")

### Android Package Signing

#### Keystore
- `ANDROID_KEYSTORE_BASE64`: Base64-encoded Android keystore file (.jks)
  ```bash
  cat android-keystore.jks | base64 -w 0
  ```
- `ANDROID_KEYSTORE_PASSWORD`: Password for the keystore
- `ANDROID_KEY_ALIAS`: Alias of the key in the keystore

### iOS Package Signing

#### Certificates and Provisioning
- `IOS_CERT_BASE64`: Base64-encoded PKCS#12 certificate file (.p12)
  ```bash
  cat ios-certificate.p12 | base64 -w 0
  ```
- `IOS_CERT_PASSWORD`: Password for the iOS certificate
- `IOS_KEYCHAIN_PASSWORD`: Password for the temporary keychain
- `IOS_PROVISIONING_PROFILE_BASE64`: Base64-encoded provisioning profile (.mobileprovision)
  ```bash
  cat TrafficSim.mobileprovision | base64 -w 0
  ```

### GitHub Package Registry

#### Authentication
- `GITHUB_TOKEN`: Automatically provided by GitHub Actions for package publishing

## Setting Up Signing Keys

### 1. Generate GPG Key (for Linux packages)
```bash
gpg --gen-key
# Follow prompts to create key
gpg --list-secret-keys --keyid-format LONG
# Note the key ID from the output
```

### 2. Windows Code Signing Certificate
- Purchase a code signing certificate from a trusted CA (DigiCert, Sectigo, etc.)
- Convert to PKCS#12 format if needed:
  ```bash
  openssl pkcs12 -export -out certificate.p12 -inkey private.key -in certificate.crt
  ```

### 3. macOS Developer Certificate
- Enroll in Apple Developer Program
- Create a Developer ID Application certificate in Apple Developer portal
- Download and export as PKCS#12 format

### 4. Android Keystore
```bash
keytool -genkey -v -keystore android-keystore.jks -alias trafficsim -keyalg RSA -keysize 2048 -validity 10000
```

### 5. iOS Certificate and Provisioning Profile
- Create App ID in Apple Developer portal
- Generate iOS Distribution certificate
- Create distribution provisioning profile
- Download both files

## Security Best Practices

1. **Use separate certificates for different environments** (dev/staging/prod)
2. **Rotate certificates before expiration**
3. **Store certificates securely** outside of version control
4. **Use minimum required permissions** for signing operations
5. **Monitor certificate usage** and set up expiration alerts
6. **Test signing process** in a staging environment first

## Troubleshooting

### Common Issues

#### GPG Signing Fails
- Ensure GPG key is not expired
- Check that the key has signing capabilities
- Verify base64 encoding is correct

#### Windows Signing Fails
- Ensure certificate is valid and not expired
- Check that signtool.exe is available in the expected path
- Verify timestamp server is accessible

#### macOS Signing Fails
- Ensure Xcode command line tools are installed
- Check that Developer ID is correctly formatted
- Verify certificate chain is complete

#### Android Signing Fails
- Ensure keystore file is not corrupted
- Check that alias exists in keystore
- Verify Gradle configuration is correct

#### iOS Signing Fails
- Ensure provisioning profile matches bundle ID
- Check that certificate is valid for distribution
- Verify Xcode project configuration

## Package Distribution

### GitHub Packages (NPM Registry)
Packages are published to: `@rhino11/trafficsim`

Install with:
```bash
npm install @rhino11/trafficsim
```

### Docker Registry
Platform-specific containers available at:
- `ghcr.io/rhino11/trafficsim-linux:latest`
- `ghcr.io/rhino11/trafficsim-windows:latest`
- `ghcr.io/rhino11/trafficsim-macos:latest`
- `ghcr.io/rhino11/trafficsim-android:latest`
- `ghcr.io/rhino11/trafficsim-ios:latest`

### GitHub Releases
All packages are attached to GitHub releases with detailed release notes.

## Maintenance

### Regular Tasks
1. **Monitor certificate expiration** (set calendar reminders)
2. **Update signing keys** when needed
3. **Test package installation** on target platforms
4. **Review security vulnerabilities** in dependencies
5. **Update build tools** (Gradle, Xcode, etc.) regularly
