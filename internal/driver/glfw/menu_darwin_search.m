//go:build darwin && !no_native_menus

#import <Foundation/Foundation.h>
#import <AppKit/AppKit.h>

// External Go callbacks
extern void menuSearchCallback(const char* searchTerm);
extern void menuItemFoundCallback(int menuItemId);

// Forward declarations
void clearDarwinMenuSearch(void);
void searchMenuItemRecursive(NSMenuItem *item, NSString *searchTerm);
void clearMenuItemHighlight(NSMenuItem *item);
NSString* getMenuPath(NSMenuItem* item);
void searchDarwinMenus(const char* searchTerm);

// Custom search field for native macOS menus
@interface FyneMenuSearchField : NSSearchField <NSSearchFieldDelegate>
@property (nonatomic, assign) int searchMenuItemTag;
@end

@implementation FyneMenuSearchField

- (instancetype)initWithFrame:(NSRect)frameRect {
    self = [super initWithFrame:frameRect];
    if (self) {
        [self setPlaceholderString:@"Search All Menus..."];
        [self setDelegate:self];  // Set delegate to self to receive text change notifications
        [[self cell] setSendsSearchStringImmediately:YES];
        [[self cell] setSendsWholeSearchString:YES];
    }
    return self;
}

// This method is called every time the text changes
- (void)controlTextDidChange:(NSNotification *)notification {
    NSString *searchString = [self stringValue];
    
    // Perform the search immediately as user types
    searchDarwinMenus([searchString UTF8String]);
    
    // Also call the Go callback
    menuSearchCallback([searchString UTF8String]);
}

// This is called when Enter is pressed or focus is lost
- (void)controlTextDidEndEditing:(NSNotification *)notification {
    // Optional: clear search on Enter or focus loss
}

@end

// Custom NSMenuItem that contains a search field
@interface FyneSearchMenuItem : NSMenuItem
@property (nonatomic, strong) FyneMenuSearchField *searchField;
@end

@implementation FyneSearchMenuItem
@synthesize searchField;

- (instancetype)init {
    self = [super init];
    if (self) {
        // Create a search field with initial frame
        NSRect frame = NSMakeRect(0, 0, 200, 22);
        self.searchField = [[FyneMenuSearchField alloc] initWithFrame:frame];
        
        // Calculate the width needed for the placeholder text
        NSString *placeholderString = [self.searchField placeholderString];
        if (!placeholderString) {
            placeholderString = @"Search All Menus...";
        }
        
        // Get the search field's font for accurate measurement
        NSFont *searchFont = [self.searchField font];
        if (!searchFont) {
            searchFont = [NSFont systemFontOfSize:[NSFont systemFontSize]];
        }
        
        NSDictionary *attributes = @{NSFontAttributeName: searchFont};
        NSSize textSize = [placeholderString sizeWithAttributes:attributes];
        
        // Add padding for UI elements:
        // Search icon: ~18px, left padding: 8px, right padding: 8px, clear button: 18px
        CGFloat requiredWidth = textSize.width + 52;
        
        // Set minimum and maximum widths
        CGFloat minWidth = 140;
        CGFloat maxWidth = 280;
        CGFloat finalWidth = MAX(minWidth, MIN(requiredWidth, maxWidth));
        
        // Update frame with calculated width
        frame.size.width = finalWidth;
        self.searchField.frame = frame;
        
        // Create a container view with the calculated width
        NSView *containerView = [[NSView alloc] initWithFrame:frame];
        [containerView addSubview:self.searchField];
        
        // Set the view as the menu item's view
        [self setView:containerView];
        
        // Make it not closable on click
        [self setEnabled:YES];
    }
    return self;
}

@end

// Function to create a search menu item
const void* createDarwinSearchMenuItem() {
    FyneSearchMenuItem *searchItem = [[FyneSearchMenuItem alloc] init];
    return (void*)searchItem;
}

// Global variables to store search state
static NSMenu *fileMenu = nil;
static NSMutableArray *foundMenuItems = nil;
static NSMutableArray *searchResultItems = nil;
static int searchFieldIndex = -1;

// Function to insert search item into a menu
void insertDarwinSearchMenuItem(const void* m, int index) {
    NSMenu* menu = (NSMenu*)m;
    fileMenu = menu;  // Store reference to the File menu
    searchFieldIndex = index;
    
    FyneSearchMenuItem *searchItem = [[FyneSearchMenuItem alloc] init];
    
    if (index < 0 || index > [menu numberOfItems]) {
        [menu addItem:searchItem];
    } else {
        [menu insertItem:searchItem atIndex:index];
    }
    
    // Add separator after search
    [menu insertItem:[NSMenuItem separatorItem] atIndex:index + 1];
    
    foundMenuItems = [[NSMutableArray alloc] init];
    searchResultItems = [[NSMutableArray alloc] init];
}

// Function to search through all menus and display results
void searchDarwinMenus(const char* searchTerm) {
    NSString *search = [NSString stringWithUTF8String:searchTerm];
    NSMenu *mainMenu = [NSApp mainMenu];
    
    // Clear previous search result items from the File menu
    for (NSMenuItem *item in searchResultItems) {
        [fileMenu removeItem:item];
    }
    [searchResultItems removeAllObjects];
    [foundMenuItems removeAllObjects];
    
    if ([search length] == 0) {
        // No search term, clear everything
        clearDarwinMenuSearch();
        return;
    }
    
    // Search through all menu items recursively
    for (NSMenuItem *item in [mainMenu itemArray]) {
        searchMenuItemRecursive(item, search);
    }
    
    // Add result items directly to the File menu
    int insertIndex = searchFieldIndex + 2;  // After search field and separator
    
    if ([foundMenuItems count] > 0) {
        // Add a header separator
        NSMenuItem *headerSeparator = [NSMenuItem separatorItem];
        [fileMenu insertItem:headerSeparator atIndex:insertIndex];
        [searchResultItems addObject:headerSeparator];
        insertIndex++;
        
        // Add header
        NSMenuItem *headerItem = [[NSMenuItem alloc] initWithTitle:[NSString stringWithFormat:@"— %ld Results —", [foundMenuItems count]] action:nil keyEquivalent:@""];
        [headerItem setEnabled:NO];
        [fileMenu insertItem:headerItem atIndex:insertIndex];
        [searchResultItems addObject:headerItem];
        insertIndex++;
        
        // Add found items to the File menu
        int maxResults = 10;  // Limit to prevent menu from getting too long
        int count = 0;
        for (NSMenuItem *foundItem in foundMenuItems) {
            if (count >= maxResults) break;
            
            NSString *menuPath = getMenuPath(foundItem);
            NSMenuItem *resultItem = [[NSMenuItem alloc] initWithTitle:menuPath action:foundItem.action keyEquivalent:@""];
            [resultItem setTarget:foundItem.target];
            [resultItem setTag:foundItem.tag];
            [resultItem setEnabled:foundItem.enabled];
            [fileMenu insertItem:resultItem atIndex:insertIndex];
            [searchResultItems addObject:resultItem];
            insertIndex++;
            count++;
        }
        
        if ([foundMenuItems count] > maxResults) {
            NSMenuItem *moreItem = [[NSMenuItem alloc] initWithTitle:[NSString stringWithFormat:@"...and %ld more", [foundMenuItems count] - maxResults] action:nil keyEquivalent:@""];
            [moreItem setEnabled:NO];
            [fileMenu insertItem:moreItem atIndex:insertIndex];
            [searchResultItems addObject:moreItem];
            insertIndex++;
        }
        
        // Add footer separator
        NSMenuItem *footerSeparator = [NSMenuItem separatorItem];
        [fileMenu insertItem:footerSeparator atIndex:insertIndex];
        [searchResultItems addObject:footerSeparator];
    } else {
        // No results found
        NSMenuItem *separatorItem = [NSMenuItem separatorItem];
        [fileMenu insertItem:separatorItem atIndex:insertIndex];
        [searchResultItems addObject:separatorItem];
        insertIndex++;
        
        NSMenuItem *noResultsItem = [[NSMenuItem alloc] initWithTitle:@"No matching menu items found" action:nil keyEquivalent:@""];
        [noResultsItem setEnabled:NO];
        [fileMenu insertItem:noResultsItem atIndex:insertIndex];
        [searchResultItems addObject:noResultsItem];
        insertIndex++;
        
        NSMenuItem *footerSeparator = [NSMenuItem separatorItem];
        [fileMenu insertItem:footerSeparator atIndex:insertIndex];
        [searchResultItems addObject:footerSeparator];
    }
    
    [fileMenu update];
}

// Helper function to get the full menu path for an item
NSString* getMenuPath(NSMenuItem* item) {
    NSMutableArray *path = [[NSMutableArray alloc] init];
    NSMenuItem *current = item;
    
    while (current != nil) {
        if ([current title] && [[current title] length] > 0 && ![current isSeparatorItem]) {
            [path insertObject:[current title] atIndex:0];
        }
        
        // Navigate up to parent
        NSMenu *parentMenu = [current menu];
        if (!parentMenu || ![parentMenu supermenu]) {
            break;
        }
        
        // Find the parent menu item
        NSMenu *superMenu = [parentMenu supermenu];
        current = nil;
        for (NSMenuItem *parentItem in [superMenu itemArray]) {
            if ([parentItem submenu] == parentMenu) {
                current = parentItem;
                break;
            }
        }
    }
    
    return [path componentsJoinedByString:@" → "];
}

void searchMenuItemRecursive(NSMenuItem *item, NSString *searchTerm) {
    // Skip separators and empty titles
    if ([item isSeparatorItem] || [[item title] length] == 0) {
        // Still search submenus
        NSMenu *submenu = [item submenu];
        if (submenu) {
            for (NSMenuItem *subItem in [submenu itemArray]) {
                searchMenuItemRecursive(subItem, searchTerm);
            }
        }
        return;
    }
    
    // Check if this item matches
    NSString *title = [item title];
    NSRange range = [title rangeOfString:searchTerm options:NSCaseInsensitiveSearch];
    
    if (range.location != NSNotFound && [item action] != nil && [item isEnabled]) {
        // This item matches - add to results
        [foundMenuItems addObject:item];
        
        // Call back to Go with the menu item's tag if it has one
        if ([item tag] >= 5000) {
            menuItemFoundCallback([item tag] - 5000);
        }
    }
    
    // Search submenu if exists
    NSMenu *submenu = [item submenu];
    if (submenu) {
        for (NSMenuItem *subItem in [submenu itemArray]) {
            searchMenuItemRecursive(subItem, searchTerm);
        }
    }
}

void clearDarwinMenuSearch() {
    // Clear any search highlighting
    NSMenu *mainMenu = [NSApp mainMenu];
    for (NSMenuItem *item in [mainMenu itemArray]) {
        clearMenuItemHighlight(item);
    }
}

void clearMenuItemHighlight(NSMenuItem *item) {
    // Reset any custom highlighting
    // This is where we'd remove any visual indicators
    
    NSMenu *submenu = [item submenu];
    if (submenu) {
        for (NSMenuItem *subItem in [submenu itemArray]) {
            clearMenuItemHighlight(subItem);
        }
    }
}

// Alternative: Use macOS's built-in Help menu search
void enableDarwinHelpMenuSearch() {
    // This enables the standard macOS Help menu search
    // which automatically indexes all menu items
    NSMenu *helpMenu = nil;
    NSMenu *mainMenu = [NSApp mainMenu];
    
    // Find the Help menu
    for (NSMenuItem *item in [mainMenu itemArray]) {
        if ([[item title] isEqualToString:@"Help"]) {
            helpMenu = [item submenu];
            break;
        }
    }
    
    if (!helpMenu) {
        // Create Help menu if it doesn't exist
        helpMenu = [[NSMenu alloc] initWithTitle:@"Help"];
        NSMenuItem *helpMenuItem = [[NSMenuItem alloc] initWithTitle:@"Help" action:nil keyEquivalent:@""];
        [helpMenuItem setSubmenu:helpMenu];
        [mainMenu addItem:helpMenuItem];
    }
    
    // macOS automatically adds search to the Help menu
    // We just need to ensure the menu items are properly set up
    [NSApp setHelpMenu:helpMenu];
}
