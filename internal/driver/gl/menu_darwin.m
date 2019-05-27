// +build !no_native_menus

#import <Foundation/Foundation.h>
#import <AppKit/AppKit.h>

extern void menu_callback(int);

@interface MenuItem : NSObject {

}
@end

@implementation MenuItem
- (void) tapped:(id) sender {
    menu_callback([sender tag]);
}
@end

NSMenu* nativeMainMenu() {
    NSApplication* app = [NSApplication sharedApplication];
    NSMenu* menu = [app mainMenu];//[[NSMenu alloc] initWithTitle:@"Fyne1"];

    return menu;
}

NSMenu* menu;

void createDarwinMenu(const char* label) {
    menu = [[NSMenu alloc] initWithTitle:[NSString stringWithUTF8String:label]];

}

void addDarwinMenuItem(const char* label, int id) {
    MenuItem* tapper = [[MenuItem alloc] init];

    NSMenuItem* item = [[NSMenuItem alloc] initWithTitle:[NSString stringWithUTF8String:label]
        action:@selector(tapped:) keyEquivalent:@""];
    [item setTarget:tapper];
    [item setTag:id];
    [menu addItem:item];
}

void completeDarwinMenu() {
    NSMenu* main = nativeMainMenu();

    NSMenuItem* top = [[NSMenuItem alloc] initWithTitle:@"" action:nil keyEquivalent:@""];
    [main addItem:top];
    [main setSubmenu:menu forItem:top];
}
