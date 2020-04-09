// +build darwin
// +build ios

#import <stdbool.h>
#import <Foundation/Foundation.h>

bool getPreferenceBool(const char* key, bool fallback) {
    NSString *nskey = [NSString stringWithUTF8String:key];
    id obj = [[NSUserDefaults standardUserDefaults] objectForKey:nskey];

    if (obj == nil) {
        return fallback;
    } else {
        return [[NSUserDefaults standardUserDefaults] boolForKey:nskey] == YES;
    }
}

void setPreferenceBool(const char* key, bool value) {
    NSString *nskey = [NSString stringWithUTF8String:key];
    if (value == false) {
        [[NSUserDefaults standardUserDefaults] setBool:NO forKey:nskey];
    } else {
        [[NSUserDefaults standardUserDefaults] setBool:YES forKey:nskey];
    }
}

float getPreferenceFloat(const char* key, float fallback) {
    NSString *nskey = [NSString stringWithUTF8String:key];
    id obj = [[NSUserDefaults standardUserDefaults] objectForKey:nskey];

    if (obj == nil) {
        return fallback;
    } else {
        return [[NSUserDefaults standardUserDefaults] floatForKey:nskey];
    }
}

void setPreferenceFloat(const char* key, float value) {
    NSString *nskey = [NSString stringWithUTF8String:key];
    [[NSUserDefaults standardUserDefaults] setFloat:value forKey:nskey];
}

int getPreferenceInt(const char* key, int fallback) {
    NSString *nskey = [NSString stringWithUTF8String:key];
    id obj = [[NSUserDefaults standardUserDefaults] objectForKey:nskey];

    if (obj == nil) {
        return fallback;
    } else {
        return [[NSUserDefaults standardUserDefaults] integerForKey:nskey];
    }
}

void setPreferenceInt(const char* key, int value) {
    NSString *nskey = [NSString stringWithUTF8String:key];
    [[NSUserDefaults standardUserDefaults] setInteger:value forKey:nskey];
}

const char* getPreferenceString(const char* key, const char* fallback) {
    NSString *nskey = [NSString stringWithUTF8String:key];
    NSString *ret = [[NSUserDefaults standardUserDefaults] stringForKey:nskey];

    if (ret == nil) {
        return fallback;
    } else {
        return [ret UTF8String];
    }
}

void setPreferenceString(const char* key, const char* value) {
    NSString *nskey = [NSString stringWithUTF8String:key];
    NSString *nsval = [NSString stringWithUTF8String:value];
    [[NSUserDefaults standardUserDefaults] setObject:nsval forKey:nskey];
}

void removePreferenceValue(const char* key) {
    NSString *nskey = [NSString stringWithUTF8String:key];
    [[NSUserDefaults standardUserDefaults] removeObjectForKey:nskey];
}
