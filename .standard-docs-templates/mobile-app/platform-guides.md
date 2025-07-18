# {{PROJECT_NAME}} - Platform Guides

## Platform Overview

**Framework:** {{MOBILE_FRAMEWORK}}
**Target Platforms:** {{TARGET_PLATFORMS}}

## iOS Development

### Requirements

- Xcode {{XCODE_VERSION}}
- iOS {{MIN_IOS_VERSION}}+
- CocoaPods {{COCOAPODS_VERSION}}

### Setup

1. Install Xcode from App Store
2. Install CocoaPods:

```bash
sudo gem install cocoapods
```

3. Install iOS dependencies:

```bash
cd ios && pod install
```

### Build & Run

```bash
# Debug build
{{IOS_DEBUG_COMMAND}}

# Release build
{{IOS_RELEASE_COMMAND}}
```

### iOS-Specific Features

{{#each IOS_FEATURES}}
#### {{name}}

**Implementation:** `{{implementation}}`
**Usage:** {{usage}}
{{/each}}

### App Store Deployment

{{#if HAS_IOS_DEPLOYMENT}}
**Bundle ID:** {{IOS_BUNDLE_ID}}
**Certificate:** {{IOS_CERTIFICATE}}
**Provisioning Profile:** {{IOS_PROVISIONING_PROFILE}}

**Build for App Store:**

```bash
{{IOS_APPSTORE_BUILD_COMMAND}}
```

{{/if}}

## Android Development

### Requirements

- Android Studio {{ANDROID_STUDIO_VERSION}}
- Android SDK {{ANDROID_SDK_VERSION}}
- Java {{JAVA_VERSION}}

### Setup

1. Install Android Studio
2. Configure Android SDK
3. Install Android dependencies:

```bash
cd android && ./gradlew build
```

### Build & Run

```bash
# Debug build
{{ANDROID_DEBUG_COMMAND}}

# Release build
{{ANDROID_RELEASE_COMMAND}}
```

### Android-Specific Features

{{#each ANDROID_FEATURES}}
#### {{name}}

**Implementation:** `{{implementation}}`
**Usage:** {{usage}}
{{/each}}

### Play Store Deployment

{{#if HAS_ANDROID_DEPLOYMENT}}
**Package Name:** {{ANDROID_PACKAGE_NAME}}
**Keystore:** {{ANDROID_KEYSTORE}}
**Build Type:** {{ANDROID_BUILD_TYPE}}

**Build for Play Store:**

```bash
{{ANDROID_PLAYSTORE_BUILD_COMMAND}}
```

{{/if}}

## Cross-Platform Features

### Navigation

{{#if HAS_NAVIGATION}}
**Library:** {{NAVIGATION_LIBRARY}}
**Configuration:** `{{NAVIGATION_CONFIG}}`

**Route Structure:**

{{#each NAVIGATION_ROUTES}}
- **{{name}}:** {{component}} - {{description}}
{{/each}}

{{/if}}

### State Management

{{#if HAS_STATE_MANAGEMENT}}
**Library:** {{STATE_MANAGEMENT_LIBRARY}}
**Store Location:** `{{STORE_LOCATION}}`

**State Structure:**

{{#each STATE_MODULES}}
- **{{name}}:** {{description}}
{{/each}}

{{/if}}

### API Integration

{{#if HAS_API_INTEGRATION}}
**Base URL:** {{API_BASE_URL}}
**Authentication:** {{API_AUTH_METHOD}}

**API Modules:**

{{#each API_MODULES}}
- **{{name}}:** `{{path}}` - {{description}}
{{/each}}

{{/if}}

### Storage

{{#if HAS_LOCAL_STORAGE}}
**Local Storage:** {{LOCAL_STORAGE_LIBRARY}}
**Secure Storage:** {{SECURE_STORAGE_LIBRARY}}

**Storage Keys:**

{{#each STORAGE_KEYS}}
- `{{key}}`: {{description}}
{{/each}}

{{/if}}

### Push Notifications

{{#if HAS_PUSH_NOTIFICATIONS}}
**Service:** {{PUSH_NOTIFICATION_SERVICE}}
**Configuration:** `{{PUSH_CONFIG}}`

**Notification Types:**

{{#each NOTIFICATION_TYPES}}
- **{{name}}:** {{description}}
{{/each}}

{{/if}}

## Testing

### Unit Tests

```bash
{{UNIT_TEST_COMMAND}}
```

### Integration Tests

```bash
{{INTEGRATION_TEST_COMMAND}}
```

### E2E Tests

{{#if HAS_E2E_TESTS}}
**Framework:** {{E2E_FRAMEWORK}}

```bash
{{E2E_TEST_COMMAND}}
```

{{/if}}

### Device Testing

{{#if HAS_DEVICE_TESTING}}
**iOS Simulator:** {{IOS_SIMULATOR_COMMAND}}
**Android Emulator:** {{ANDROID_EMULATOR_COMMAND}}

{{/if}}

## Performance Optimization

### Bundle Analysis

{{#if HAS_BUNDLE_ANALYSIS}}

```bash
{{BUNDLE_ANALYSIS_COMMAND}}
```

{{/if}}

### Memory Management

{{#each MEMORY_OPTIMIZATIONS}}
- **{{area}}:** {{technique}}
{{/each}}

### Network Optimization

{{#each NETWORK_OPTIMIZATIONS}}
- **{{area}}:** {{technique}}
{{/each}}

## Platform-Specific Debugging

### iOS Debugging

{{#if HAS_IOS_DEBUGGING}}
**Xcode Debugger:** {{IOS_DEBUGGER_SETUP}}
**Instruments:** {{IOS_INSTRUMENTS_SETUP}}

{{/if}}

### Android Debugging

{{#if HAS_ANDROID_DEBUGGING}}
**Android Studio Debugger:** {{ANDROID_DEBUGGER_SETUP}}
**ADB Commands:** {{ADB_COMMANDS}}

{{/if}}

## Deployment

### iOS Deployment

{{#if HAS_IOS_DEPLOYMENT}}
**TestFlight:** {{TESTFLIGHT_SETUP}}
**App Store:** {{APPSTORE_SETUP}}

{{/if}}

### Android Deployment

{{#if HAS_ANDROID_DEPLOYMENT}}
**Internal Testing:** {{INTERNAL_TESTING_SETUP}}
**Play Store:** {{PLAYSTORE_SETUP}}

{{/if}}

---
*Last updated: {{TIMESTAMP}}*
*Generated with [standard-docs](https://github.com/johnplummer/standard-docs)*

