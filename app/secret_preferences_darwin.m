//go:build !ci && !wasm && !test_web_driver && !mobile && !tinygo

#import <Foundation/Foundation.h>
#import <Security/Security.h>

#include <stdlib.h>
#include <string.h>

static NSMutableDictionary *keychainQuery(const char *service, const char *account) {
    NSMutableDictionary *query = [NSMutableDictionary dictionary];
    query[(__bridge id)kSecClass] = (__bridge id)kSecClassGenericPassword;
    query[(__bridge id)kSecAttrService] = [NSString stringWithUTF8String:service];
    query[(__bridge id)kSecAttrAccount] = [NSString stringWithUTF8String:account];
    return query;
}

// keychainLoad reads the item into a malloc'd buffer that the caller must free.
// Returns 0 on success, 1 if no item exists, or the OSStatus error otherwise.
int keychainLoad(const char *service, const char *account, void **data, int *length) {
    @autoreleasepool {
        NSMutableDictionary *query = keychainQuery(service, account);
        query[(__bridge id)kSecReturnData] = @YES;
        query[(__bridge id)kSecMatchLimit] = (__bridge id)kSecMatchLimitOne;

        CFTypeRef result = NULL;
        OSStatus status = SecItemCopyMatching((__bridge CFDictionaryRef)query, &result);
        if (status == errSecItemNotFound) {
            return 1;
        }
        if (status != errSecSuccess) {
            return (int)status;
        }

        NSData *found = (__bridge NSData *)result;
        *length = (int)[found length];
        *data = malloc(*length > 0 ? *length : 1);
        memcpy(*data, [found bytes], *length);
        CFRelease(result);
        return 0;
    }
}

// keychainSave updates the item, creating it if required.
// Returns 0 on success or the OSStatus error otherwise.
int keychainSave(const char *service, const char *account, const void *data, int length) {
    @autoreleasepool {
        NSData *value = [NSData dataWithBytes:data length:length];
        NSMutableDictionary *query = keychainQuery(service, account);

        NSDictionary *update = @{(__bridge id)kSecValueData: value};
        OSStatus status = SecItemUpdate((__bridge CFDictionaryRef)query, (__bridge CFDictionaryRef)update);
        if (status != errSecItemNotFound) {
            return status == errSecSuccess ? 0 : (int)status;
        }

        query[(__bridge id)kSecValueData] = value;
#if TARGET_OS_IPHONE
        // allow background saves once the device has been unlocked since boot
        query[(__bridge id)kSecAttrAccessible] = (__bridge id)kSecAttrAccessibleAfterFirstUnlockThisDeviceOnly;
#endif
        status = SecItemAdd((__bridge CFDictionaryRef)query, NULL);
        return status == errSecSuccess ? 0 : (int)status;
    }
}
