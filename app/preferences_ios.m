// +build darwin
// +build ios

#import <Foundation/Foundation.h>

int getPreferenceBool(const char* key, int fallback) {
    NSString *nskey = [NSString stringWithUTF8String:key];
    id obj = [[NSUserDefaults standardUserDefaults] objectForKey:nskey];

    if (obj == nil) {
        return fallback;
    } else {
        if ([[NSUserDefaults standardUserDefaults] boolForKey:nskey]) {
            return 1;
        } else {
            return 0;
        }
    }
}

void setPreferenceBool(const char* key, int value) {
    NSString *nskey = [NSString stringWithUTF8String:key];
    if (value == 0) {
        [[NSUserDefaults standardUserDefaults] setBool:FALSE forKey:nskey];
    } else {
        [[NSUserDefaults standardUserDefaults] setBool:TRUE forKey:nskey];
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
