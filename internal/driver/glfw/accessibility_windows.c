//go:build accessibility && windows

#define CINTERFACE
#define COBJMACROS
#include "accessibility_windows.h"
#include <windows.h>
#include <ole2.h>
#include <oleacc.h>

// ============================================================
// UIA type definitions (manual, for MinGW/CGo compatibility)
// ============================================================

typedef struct FyneUIAElement FyneUIAElement;

typedef struct IRawSimple { struct IRawSimpleVtbl* lpVtbl; } IRawSimple;
typedef struct IRawFragment { struct IRawFragmentVtbl* lpVtbl; } IRawFragment;
typedef struct IRawFragRoot { struct IRawFragRootVtbl* lpVtbl; } IRawFragRoot;

typedef int PROPERTYID;
typedef int PATTERNID;
typedef int EVENTID;

enum UIAProviderOptions {
    UIAProviderOptions_ServerSideProvider = 0x2,
    UIAProviderOptions_UseComThreading = 0x20,
};

enum UIANavigateDirection {
    UIANavigateDirection_Parent = 0,
    UIANavigateDirection_NextSibling = 1,
    UIANavigateDirection_PreviousSibling = 2,
    UIANavigateDirection_FirstChild = 3,
    UIANavigateDirection_LastChild = 4,
};

enum UIAStructureChangeType {
    UIAStructureChangeType_ChildrenInvalidated = 2,
};

typedef struct { double left, top, width, height; } UIARect;

struct IRawSimpleVtbl {
    HRESULT (STDMETHODCALLTYPE *QueryInterface)(IRawSimple*, REFIID, void**);
    ULONG   (STDMETHODCALLTYPE *AddRef)(IRawSimple*);
    ULONG   (STDMETHODCALLTYPE *Release)(IRawSimple*);
    HRESULT (STDMETHODCALLTYPE *get_ProviderOptions)(IRawSimple*, int*);
    HRESULT (STDMETHODCALLTYPE *GetPatternProvider)(IRawSimple*, PATTERNID, IUnknown**);
    HRESULT (STDMETHODCALLTYPE *GetPropertyValue)(IRawSimple*, PROPERTYID, VARIANT*);
    HRESULT (STDMETHODCALLTYPE *get_HostRawElementProvider)(IRawSimple*, IRawSimple**);
};

struct IRawFragmentVtbl {
    HRESULT (STDMETHODCALLTYPE *QueryInterface)(IRawFragment*, REFIID, void**);
    ULONG   (STDMETHODCALLTYPE *AddRef)(IRawFragment*);
    ULONG   (STDMETHODCALLTYPE *Release)(IRawFragment*);
    HRESULT (STDMETHODCALLTYPE *Navigate)(IRawFragment*, int, IRawFragment**);
    HRESULT (STDMETHODCALLTYPE *GetRuntimeId)(IRawFragment*, SAFEARRAY**);
    HRESULT (STDMETHODCALLTYPE *get_BoundingRectangle)(IRawFragment*, UIARect*);
    HRESULT (STDMETHODCALLTYPE *GetEmbeddedFragmentRoots)(IRawFragment*, SAFEARRAY**);
    HRESULT (STDMETHODCALLTYPE *SetFocus)(IRawFragment*);
    HRESULT (STDMETHODCALLTYPE *get_FragmentRoot)(IRawFragment*, IRawFragRoot**);
};

struct IRawFragRootVtbl {
    HRESULT (STDMETHODCALLTYPE *QueryInterface)(IRawFragRoot*, REFIID, void**);
    ULONG   (STDMETHODCALLTYPE *AddRef)(IRawFragRoot*);
    ULONG   (STDMETHODCALLTYPE *Release)(IRawFragRoot*);
    HRESULT (STDMETHODCALLTYPE *ElementProviderFromPoint)(IRawFragRoot*, double, double, IRawSimple**);
    HRESULT (STDMETHODCALLTYPE *GetFocus)(IRawFragRoot*, IRawFragment**);
};

// ============================================================
// Element struct
// ============================================================

struct FyneUIAElement {
    IRawSimple    simple;
    IRawFragment  fragment;
    IRawFragRoot  fragRoot;

    LONG refCount;
    int  isRoot;
    HWND hwnd;
    int  uniqueId;

    FyneUIAElement* parent;
    WCHAR* name;
    int    controlType;
    double x, y, width, height;
    int    childIndex;

    FyneUIAElement** children;
    int childCount;
    int childCapacity;
};

#define ELEM_FROM_SIMPLE(p)   ((FyneUIAElement*)(p))
#define ELEM_FROM_FRAGMENT(p) ((FyneUIAElement*)((char*)(p) - offsetof(FyneUIAElement, fragment)))
#define ELEM_FROM_FRAGROOT(p) ((FyneUIAElement*)((char*)(p) - offsetof(FyneUIAElement, fragRoot)))

// ============================================================
// GUIDs
// ============================================================

static const IID LOCAL_IID_IUnknown =
    {0x00000000,0x0000,0x0000,{0xC0,0x00,0x00,0x00,0x00,0x00,0x00,0x46}};
static const IID IID_IRawSimple =
    {0xd6dd68d1,0x86fd,0x4332,{0x86,0x66,0x9a,0xbe,0xde,0xa2,0xd2,0x4c}};
static const IID IID_IRawFragment =
    {0xf7063da8,0x8359,0x439c,{0x92,0x97,0xbb,0xc5,0x29,0x9a,0x7d,0x87}};
static const IID IID_IRawFragRoot =
    {0x620ce2a5,0xab8f,0x40a9,{0x86,0xcb,0xde,0x3c,0x75,0x59,0x9b,0x58}};

// UIA property IDs
#define UIA_ControlTypePropertyId          30003
#define UIA_NamePropertyId                 30005
#define UIA_HasKeyboardFocusPropertyId     30008
#define UIA_IsKeyboardFocusablePropertyId  30009
#define UIA_IsEnabledPropertyId            30010
#define UIA_AutomationIdPropertyId         30011
#define UIA_IsControlElementPropertyId     30016
#define UIA_IsContentElementPropertyId     30017
#define UIA_ProviderDescriptionPropertyId  30107

// UIA control type IDs
#define UIA_ButtonControlTypeId    50000
#define UIA_CheckBoxControlTypeId  50002
#define UIA_ComboBoxControlTypeId  50003
#define UIA_EditControlTypeId      50004
#define UIA_HyperlinkControlTypeId 50005
#define UIA_ImageControlTypeId     50006
#define UIA_ListItemControlTypeId  50007
#define UIA_ListControlTypeId      50008
#define UIA_MenuControlTypeId      50009
#define UIA_ProgressBarControlTypeId 50012
#define UIA_RadioButtonControlTypeId 50013
#define UIA_ScrollBarControlTypeId 50014
#define UIA_SliderControlTypeId    50015
#define UIA_SpinnerControlTypeId   50016
#define UIA_TabControlTypeId       50018
#define UIA_TabItemControlTypeId   50019
#define UIA_TextControlTypeId      50020
#define UIA_TreeControlTypeId      50023
#define UIA_TreeItemControlTypeId  50024
#define UIA_GroupControlTypeId     50026
#define UIA_HeaderControlTypeId    50034
#define UIA_DataGridControlTypeId  50028
#define UIA_PaneControlTypeId      50033
#define UIA_SeparatorControlTypeId 50038

#define UiaAppendRuntimeId  3
#define UiaRootObjectId    (-25)

#define UIA_AutomationFocusChangedEventId 20005

#define WM_FYNE_RAISE_FOCUS (WM_APP + 100)
#define WM_FYNE_FOCUS_CHILD (WM_APP + 101)

// ============================================================
// Dynamically loaded UIA functions
// ============================================================

typedef LRESULT (WINAPI *PFN_UiaReturnRawElementProvider)(HWND, WPARAM, LPARAM, void*);
typedef HRESULT (WINAPI *PFN_UiaHostProviderFromHwnd)(HWND, void**);
typedef HRESULT (WINAPI *PFN_UiaRaiseAutomationEvent)(void*, EVENTID);
typedef HRESULT (WINAPI *PFN_UiaRaiseStructureChangedEvent)(void*, int, int*, int);
typedef HRESULT (WINAPI *PFN_UiaDisconnectProvider)(void*);

static PFN_UiaReturnRawElementProvider pfnUiaReturn = NULL;
static PFN_UiaHostProviderFromHwnd     pfnUiaHost   = NULL;
static PFN_UiaRaiseAutomationEvent     pfnUiaRaiseEvent = NULL;
static PFN_UiaRaiseStructureChangedEvent pfnUiaRaiseStructure = NULL;
static PFN_UiaDisconnectProvider       pfnUiaDisconnect = NULL;
static HMODULE hUiaCore = NULL;

static void loadUiaFunctions(void) {
    if (hUiaCore) return;
    hUiaCore = LoadLibraryW(L"uiautomationcore.dll");
    if (!hUiaCore) return;
    pfnUiaReturn = (PFN_UiaReturnRawElementProvider)GetProcAddress(hUiaCore, "UiaReturnRawElementProvider");
    pfnUiaHost   = (PFN_UiaHostProviderFromHwnd)GetProcAddress(hUiaCore, "UiaHostProviderFromHwnd");
    pfnUiaRaiseEvent = (PFN_UiaRaiseAutomationEvent)GetProcAddress(hUiaCore, "UiaRaiseAutomationEvent");
    pfnUiaRaiseStructure = (PFN_UiaRaiseStructureChangedEvent)GetProcAddress(hUiaCore, "UiaRaiseStructureChangedEvent");
    pfnUiaDisconnect = (PFN_UiaDisconnectProvider)GetProcAddress(hUiaCore, "UiaDisconnectProvider");
}

// ============================================================
// Globals
// ============================================================

static FyneUIAElement* g_root = NULL;
static HWND g_hwnd = NULL;
static WNDPROC g_origWndProc = NULL;
static struct IRawSimpleVtbl   g_simpleVtbl;
static struct IRawFragmentVtbl g_fragmentVtbl;
static struct IRawFragRootVtbl g_fragRootVtbl;
static int g_vtblInit = 0;
static int g_nextId = 1;
static int g_focusedIndex = -1;

static FyneUIAElement** g_staging = NULL;
static int g_stagingCount = 0;
static int g_stagingCapacity = 0;

// ============================================================
// Helpers
// ============================================================

static WCHAR* utf8ToWide(const char* utf8) {
    if (!utf8 || !utf8[0]) {
        WCHAR* e = (WCHAR*)malloc(sizeof(WCHAR));
        if (e) e[0] = L'\0';
        return e;
    }
    int len = MultiByteToWideChar(CP_UTF8, 0, utf8, -1, NULL, 0);
    WCHAR* w = (WCHAR*)malloc(len * sizeof(WCHAR));
    if (w) MultiByteToWideChar(CP_UTF8, 0, utf8, -1, w, len);
    return w;
}

static int roleToUIA(WinAccessibilityRole role) {
    switch (role) {
    case WinAccessibilityRoleButton:      return UIA_ButtonControlTypeId;
    case WinAccessibilityRoleText:        return UIA_TextControlTypeId;
    case WinAccessibilityRoleLink:        return UIA_HyperlinkControlTypeId;
    case WinAccessibilityRoleGroup:       return UIA_GroupControlTypeId;
    case WinAccessibilityRoleCheckbox:    return UIA_CheckBoxControlTypeId;
    case WinAccessibilityRoleRadio:       return UIA_RadioButtonControlTypeId;
    case WinAccessibilityRoleSlider:      return UIA_SliderControlTypeId;
    case WinAccessibilityRoleProgressBar: return UIA_ProgressBarControlTypeId;
    case WinAccessibilityRoleTextField:   return UIA_EditControlTypeId;
    case WinAccessibilityRoleList:        return UIA_ListControlTypeId;
    case WinAccessibilityRoleListItem:    return UIA_ListItemControlTypeId;
    case WinAccessibilityRoleTree:        return UIA_TreeControlTypeId;
    case WinAccessibilityRoleTreeItem:    return UIA_TreeItemControlTypeId;
    case WinAccessibilityRoleTable:       return UIA_DataGridControlTypeId;
    case WinAccessibilityRoleTab:         return UIA_TabItemControlTypeId;
    case WinAccessibilityRoleTabList:     return UIA_TabControlTypeId;
    case WinAccessibilityRoleImage:       return UIA_ImageControlTypeId;
    case WinAccessibilityRoleHeading:     return UIA_HeaderControlTypeId;
    case WinAccessibilityRoleSeparator:   return UIA_SeparatorControlTypeId;
    default:                              return UIA_PaneControlTypeId;
    }
}

static HRESULT elemQI(FyneUIAElement* elem, REFIID riid, void** ppv) {
    if (!ppv) return E_POINTER;
    if (IsEqualIID(riid, &LOCAL_IID_IUnknown) || IsEqualIID(riid, &IID_IRawSimple)) {
        *ppv = &elem->simple;
        InterlockedIncrement(&elem->refCount);
        return S_OK;
    }
    if (IsEqualIID(riid, &IID_IRawFragment)) {
        *ppv = &elem->fragment;
        InterlockedIncrement(&elem->refCount);
        return S_OK;
    }
    if (elem->isRoot && IsEqualIID(riid, &IID_IRawFragRoot)) {
        *ppv = &elem->fragRoot;
        InterlockedIncrement(&elem->refCount);
        return S_OK;
    }
    *ppv = NULL;
    return E_NOINTERFACE;
}

// ============================================================
// IRawElementProviderSimple
// ============================================================

static HRESULT STDMETHODCALLTYPE S_QI(IRawSimple* This, REFIID riid, void** ppv) {
    return elemQI(ELEM_FROM_SIMPLE(This), riid, ppv);
}
static ULONG STDMETHODCALLTYPE S_AddRef(IRawSimple* This) {
    return InterlockedIncrement(&ELEM_FROM_SIMPLE(This)->refCount);
}
static ULONG STDMETHODCALLTYPE S_Release(IRawSimple* This) {
    FyneUIAElement* e = ELEM_FROM_SIMPLE(This);
    ULONG c = InterlockedDecrement(&e->refCount);
    if (c == 0 && !e->isRoot) { free(e->name); free(e); }
    return c;
}

static HRESULT STDMETHODCALLTYPE S_get_ProviderOptions(IRawSimple* This, int* pRetVal) {
    if (!pRetVal) return E_POINTER;
    *pRetVal = UIAProviderOptions_ServerSideProvider | UIAProviderOptions_UseComThreading;
    return S_OK;
}

static HRESULT STDMETHODCALLTYPE S_GetPatternProvider(IRawSimple* This, PATTERNID id, IUnknown** pRetVal) {
    if (!pRetVal) return E_POINTER;
    *pRetVal = NULL;
    return S_OK;
}

static HRESULT STDMETHODCALLTYPE S_GetPropertyValue(IRawSimple* This, PROPERTYID pid, VARIANT* pRetVal) {
    if (!pRetVal) return E_POINTER;
    VariantInit(pRetVal);
    FyneUIAElement* e = ELEM_FROM_SIMPLE(This);

    // Root element: let host provider handle most properties,
    // override focus-related ones (host returns IsKeyboardFocusable=False)
    if (e->isRoot) {
        switch (pid) {
        case UIA_IsKeyboardFocusablePropertyId:
            pRetVal->vt = VT_BOOL;
            pRetVal->boolVal = VARIANT_TRUE;
            break;
        case UIA_HasKeyboardFocusPropertyId:
            pRetVal->vt = VT_BOOL;
            pRetVal->boolVal = (GetForegroundWindow() == e->hwnd) ? VARIANT_TRUE : VARIANT_FALSE;
            break;
        case UIA_IsControlElementPropertyId:
        case UIA_IsContentElementPropertyId:
        case UIA_IsEnabledPropertyId:
            pRetVal->vt = VT_BOOL;
            pRetVal->boolVal = VARIANT_TRUE;
            break;
        case UIA_ProviderDescriptionPropertyId:
            pRetVal->vt = VT_BSTR;
            pRetVal->bstrVal = SysAllocString(L"Fyne Accessibility Provider");
            break;
        }
        return S_OK;
    }

    // Child element properties
    switch (pid) {
    case UIA_ControlTypePropertyId:
        pRetVal->vt = VT_I4;
        pRetVal->lVal = e->controlType;
        break;
    case UIA_NamePropertyId:
        pRetVal->vt = VT_BSTR;
        pRetVal->bstrVal = SysAllocString(e->name);
        break;
    case UIA_AutomationIdPropertyId: {
        WCHAR buf[32];
        wsprintfW(buf, L"fyne_%d", e->uniqueId);
        pRetVal->vt = VT_BSTR;
        pRetVal->bstrVal = SysAllocString(buf);
        break;
    }
    case UIA_IsControlElementPropertyId:
    case UIA_IsContentElementPropertyId:
    case UIA_IsEnabledPropertyId:
    case UIA_IsKeyboardFocusablePropertyId:
        pRetVal->vt = VT_BOOL;
        pRetVal->boolVal = VARIANT_TRUE;
        break;
    case UIA_HasKeyboardFocusPropertyId:
        pRetVal->vt = VT_BOOL;
        pRetVal->boolVal = (g_focusedIndex == e->childIndex) ? VARIANT_TRUE : VARIANT_FALSE;
        break;
    case UIA_ProviderDescriptionPropertyId:
        pRetVal->vt = VT_BSTR;
        pRetVal->bstrVal = SysAllocString(L"Fyne Accessibility Provider");
        break;
    }
    return S_OK;
}

static HRESULT STDMETHODCALLTYPE S_get_HostRawElementProvider(IRawSimple* This, IRawSimple** pRetVal) {
    if (!pRetVal) return E_POINTER;
    *pRetVal = NULL;
    FyneUIAElement* e = ELEM_FROM_SIMPLE(This);
    if (e->isRoot && pfnUiaHost) {
        pfnUiaHost(e->hwnd, (void**)pRetVal);
    }
    return S_OK;
}

// ============================================================
// IRawElementProviderFragment
// ============================================================

static HRESULT STDMETHODCALLTYPE F_QI(IRawFragment* This, REFIID riid, void** ppv) {
    return elemQI(ELEM_FROM_FRAGMENT(This), riid, ppv);
}
static ULONG STDMETHODCALLTYPE F_AddRef(IRawFragment* This) {
    return InterlockedIncrement(&ELEM_FROM_FRAGMENT(This)->refCount);
}
static ULONG STDMETHODCALLTYPE F_Release(IRawFragment* This) {
    FyneUIAElement* e = ELEM_FROM_FRAGMENT(This);
    ULONG c = InterlockedDecrement(&e->refCount);
    if (c == 0 && !e->isRoot) { free(e->name); free(e); }
    return c;
}

static HRESULT STDMETHODCALLTYPE F_Navigate(IRawFragment* This, int direction, IRawFragment** pRetVal) {
    if (!pRetVal) return E_POINTER;
    *pRetVal = NULL;
    FyneUIAElement* e = ELEM_FROM_FRAGMENT(This);
    FyneUIAElement* target = NULL;

    if (e->isRoot) {
        switch (direction) {
        case UIANavigateDirection_FirstChild:
            if (e->childCount > 0) target = e->children[0];
            break;
        case UIANavigateDirection_LastChild:
            if (e->childCount > 0) target = e->children[e->childCount - 1];
            break;
        }
    } else {
        FyneUIAElement* p = e->parent;
        switch (direction) {
        case UIANavigateDirection_Parent:
            target = p;
            break;
        case UIANavigateDirection_NextSibling:
            if (p && e->childIndex + 1 < p->childCount)
                target = p->children[e->childIndex + 1];
            break;
        case UIANavigateDirection_PreviousSibling:
            if (p && e->childIndex > 0)
                target = p->children[e->childIndex - 1];
            break;
        }
    }

    if (target) {
        *pRetVal = &target->fragment;
        InterlockedIncrement(&target->refCount);
    }
    return S_OK;
}

static HRESULT STDMETHODCALLTYPE F_GetRuntimeId(IRawFragment* This, SAFEARRAY** pRetVal) {
    if (!pRetVal) return E_POINTER;
    FyneUIAElement* e = ELEM_FROM_FRAGMENT(This);

    SAFEARRAYBOUND bound = {2, 0};
    SAFEARRAY* psa = SafeArrayCreate(VT_I4, 1, &bound);
    if (!psa) return E_OUTOFMEMORY;

    long idx = 0;
    int val = UiaAppendRuntimeId;
    SafeArrayPutElement(psa, &idx, &val);
    idx = 1;
    val = e->uniqueId;
    SafeArrayPutElement(psa, &idx, &val);

    *pRetVal = psa;
    return S_OK;
}

static HRESULT STDMETHODCALLTYPE F_get_BoundingRectangle(IRawFragment* This, UIARect* pRetVal) {
    if (!pRetVal) return E_POINTER;
    FyneUIAElement* e = ELEM_FROM_FRAGMENT(This);

    if (e->isRoot) {
        RECT rc;
        GetClientRect(e->hwnd, &rc);
        POINT pt = {0, 0};
        ClientToScreen(e->hwnd, &pt);
        pRetVal->left   = pt.x;
        pRetVal->top    = pt.y;
        pRetVal->width  = rc.right - rc.left;
        pRetVal->height = rc.bottom - rc.top;
    } else {
        POINT pt = {(LONG)e->x, (LONG)e->y};
        ClientToScreen(e->hwnd, &pt);
        pRetVal->left   = pt.x;
        pRetVal->top    = pt.y;
        pRetVal->width  = e->width;
        pRetVal->height = e->height;
    }
    return S_OK;
}

static HRESULT STDMETHODCALLTYPE F_GetEmbeddedFragmentRoots(IRawFragment* This, SAFEARRAY** pRetVal) {
    if (pRetVal) *pRetVal = NULL;
    return S_OK;
}

static HRESULT STDMETHODCALLTYPE F_SetFocus(IRawFragment* This) {
    FyneUIAElement* e = ELEM_FROM_FRAGMENT(This);
    if (!e->isRoot) {
        g_focusedIndex = e->childIndex;
        if (pfnUiaRaiseEvent) {
            pfnUiaRaiseEvent(&e->simple, UIA_AutomationFocusChangedEventId);
        }
    }
    return S_OK;
}

static HRESULT STDMETHODCALLTYPE F_get_FragmentRoot(IRawFragment* This, IRawFragRoot** pRetVal) {
    if (!pRetVal) return E_POINTER;
    *pRetVal = NULL;
    FyneUIAElement* e = ELEM_FROM_FRAGMENT(This);
    FyneUIAElement* root = e->isRoot ? e : e->parent;
    if (root) {
        *pRetVal = &root->fragRoot;
        InterlockedIncrement(&root->refCount);
    }
    return S_OK;
}

// ============================================================
// IRawElementProviderFragmentRoot
// ============================================================

static HRESULT STDMETHODCALLTYPE FR_QI(IRawFragRoot* This, REFIID riid, void** ppv) {
    return elemQI(ELEM_FROM_FRAGROOT(This), riid, ppv);
}
static ULONG STDMETHODCALLTYPE FR_AddRef(IRawFragRoot* This) {
    return InterlockedIncrement(&ELEM_FROM_FRAGROOT(This)->refCount);
}
static ULONG STDMETHODCALLTYPE FR_Release(IRawFragRoot* This) {
    return InterlockedDecrement(&ELEM_FROM_FRAGROOT(This)->refCount);
}

static HRESULT STDMETHODCALLTYPE FR_ElementProviderFromPoint(IRawFragRoot* This,
    double x, double y, IRawSimple** pRetVal) {
    if (!pRetVal) return E_POINTER;
    *pRetVal = NULL;
    FyneUIAElement* e = ELEM_FROM_FRAGROOT(This);

    POINT pt = {(LONG)x, (LONG)y};
    ScreenToClient(e->hwnd, &pt);

    for (int i = 0; i < e->childCount; i++) {
        FyneUIAElement* c = e->children[i];
        if (pt.x >= c->x && pt.x < c->x + c->width &&
            pt.y >= c->y && pt.y < c->y + c->height) {
            *pRetVal = &c->simple;
            InterlockedIncrement(&c->refCount);
            return S_OK;
        }
    }
    return S_OK;
}

static HRESULT STDMETHODCALLTYPE FR_GetFocus(IRawFragRoot* This, IRawFragment** pRetVal) {
    if (!pRetVal) return E_POINTER;
    *pRetVal = NULL;
    FyneUIAElement* e = ELEM_FROM_FRAGROOT(This);
    if (g_focusedIndex >= 0 && g_focusedIndex < e->childCount) {
        FyneUIAElement* child = e->children[g_focusedIndex];
        *pRetVal = &child->fragment;
        InterlockedIncrement(&child->refCount);
    }
    return S_OK;
}

// ============================================================
// Vtable setup
// ============================================================

static void initVtbls(void) {
    if (g_vtblInit) return;

    g_simpleVtbl.QueryInterface = S_QI;
    g_simpleVtbl.AddRef = S_AddRef;
    g_simpleVtbl.Release = S_Release;
    g_simpleVtbl.get_ProviderOptions = S_get_ProviderOptions;
    g_simpleVtbl.GetPatternProvider = S_GetPatternProvider;
    g_simpleVtbl.GetPropertyValue = S_GetPropertyValue;
    g_simpleVtbl.get_HostRawElementProvider = S_get_HostRawElementProvider;

    g_fragmentVtbl.QueryInterface = F_QI;
    g_fragmentVtbl.AddRef = F_AddRef;
    g_fragmentVtbl.Release = F_Release;
    g_fragmentVtbl.Navigate = F_Navigate;
    g_fragmentVtbl.GetRuntimeId = F_GetRuntimeId;
    g_fragmentVtbl.get_BoundingRectangle = F_get_BoundingRectangle;
    g_fragmentVtbl.GetEmbeddedFragmentRoots = F_GetEmbeddedFragmentRoots;
    g_fragmentVtbl.SetFocus = F_SetFocus;
    g_fragmentVtbl.get_FragmentRoot = F_get_FragmentRoot;

    g_fragRootVtbl.QueryInterface = FR_QI;
    g_fragRootVtbl.AddRef = FR_AddRef;
    g_fragRootVtbl.Release = FR_Release;
    g_fragRootVtbl.ElementProviderFromPoint = FR_ElementProviderFromPoint;
    g_fragRootVtbl.GetFocus = FR_GetFocus;

    g_vtblInit = 1;
}

static FyneUIAElement* createElement(int isRoot, HWND hwnd) {
    FyneUIAElement* e = (FyneUIAElement*)calloc(1, sizeof(FyneUIAElement));
    if (!e) return NULL;
    e->simple.lpVtbl   = &g_simpleVtbl;
    e->fragment.lpVtbl  = &g_fragmentVtbl;
    e->fragRoot.lpVtbl  = &g_fragRootVtbl;
    e->refCount = 1;
    e->isRoot = isRoot;
    e->hwnd = hwnd;
    e->uniqueId = g_nextId++;
    return e;
}

// ============================================================
// Window subclass
// ============================================================

static void focusChild(int index) {
    if (!g_root || index < 0 || index >= g_root->childCount) return;
    if (index == g_focusedIndex) return;
    PostMessageW(g_hwnd, WM_FYNE_FOCUS_CHILD, (WPARAM)index, 0);
}

static int hitTestChild(int clientX, int clientY) {
    if (!g_root) return -1;
    for (int i = 0; i < g_root->childCount; i++) {
        FyneUIAElement* c = g_root->children[i];
        if (clientX >= c->x && clientX < c->x + c->width &&
            clientY >= c->y && clientY < c->y + c->height) {
            return i;
        }
    }
    return -1;
}

static LRESULT CALLBACK AccessibilityWndProc(HWND hwnd, UINT msg, WPARAM wParam, LPARAM lParam) {
    if (msg == WM_GETOBJECT) {
        if (lParam == (LPARAM)UiaRootObjectId && g_root && pfnUiaReturn) {
            return pfnUiaReturn(hwnd, wParam, lParam, &g_root->simple);
        }
        return DefWindowProcW(hwnd, msg, wParam, lParam);
    }

    if (msg == WM_SETFOCUS || (msg == WM_ACTIVATE && LOWORD(wParam) != 0)) {
        PostMessageW(hwnd, WM_FYNE_RAISE_FOCUS, 0, 0);
    }

    if (msg == WM_FYNE_RAISE_FOCUS) {
        if (g_root && pfnUiaRaiseEvent) {
            pfnUiaRaiseEvent(&g_root->simple, UIA_AutomationFocusChangedEventId);
        }
        return 0;
    }

    if (msg == WM_FYNE_FOCUS_CHILD) {
        int index = (int)wParam;
        if (g_root && index >= 0 && index < g_root->childCount) {
            g_focusedIndex = index;
            if (pfnUiaRaiseEvent) {
                pfnUiaRaiseEvent(&g_root->children[index]->simple, UIA_AutomationFocusChangedEventId);
            }
        }
        return 0;
    }

    if (msg == WM_LBUTTONDOWN) {
        int hit = hitTestChild((int)(short)LOWORD(lParam), (int)(short)HIWORD(lParam));
        if (hit >= 0) {
            focusChild(hit);
        }
    }

    if (msg == WM_KEYDOWN && wParam == VK_TAB) {
        if (g_root && g_root->childCount > 0) {
            int next;
            if (GetKeyState(VK_SHIFT) & 0x8000) {
                next = (g_focusedIndex <= 0) ? g_root->childCount - 1 : g_focusedIndex - 1;
            } else {
                next = (g_focusedIndex + 1) % g_root->childCount;
            }
            focusChild(next);
        }
    }

    return CallWindowProcW(g_origWndProc, hwnd, msg, wParam, lParam);
}

// ============================================================
// Public API
// ============================================================

void WinAccessibilitySetWindow(void* hwnd) {
    HWND h = (HWND)hwnd;
    if (h == g_hwnd && g_root) return;

    loadUiaFunctions();

    if (g_hwnd && g_origWndProc) {
        SetWindowLongPtrW(g_hwnd, GWLP_WNDPROC, (LONG_PTR)g_origWndProc);
        g_origWndProc = NULL;
    }

    g_hwnd = h;
    initVtbls();

    if (!g_root) {
        g_root = createElement(1, h);
        if (!g_root) return;
    }
    g_root->hwnd = h;

    g_origWndProc = (WNDPROC)SetWindowLongPtrW(h, GWLP_WNDPROC, (LONG_PTR)AccessibilityWndProc);
}

void WinAccessibilityAddElement(const char* name, WinAccessibilityRole role,
    double x, double y, double width, double height) {
    if (!g_root) return;

    if (g_stagingCount >= g_stagingCapacity) {
        int newCap = g_stagingCapacity == 0 ? 16 : g_stagingCapacity * 2;
        FyneUIAElement** a = (FyneUIAElement**)realloc(g_staging, newCap * sizeof(FyneUIAElement*));
        if (!a) return;
        g_staging = a;
        g_stagingCapacity = newCap;
    }

    FyneUIAElement* child = createElement(0, g_root->hwnd);
    if (!child) return;
    child->parent = g_root;
    child->name = utf8ToWide(name);
    child->controlType = roleToUIA(role);
    child->x = x;
    child->y = y;
    child->width = width;
    child->height = height;
    child->childIndex = g_stagingCount;

    g_staging[g_stagingCount] = child;
    g_stagingCount++;
}

void WinAccessibilityClearElements(void) {
    g_stagingCount = 0;
}

void WinAccessibilityUpdate(void) {
    if (!g_root || !g_hwnd) return;

    // Check if tree structure changed (count or names/roles differ)
    int structureChanged = 0;
    if (g_stagingCount != g_root->childCount) {
        structureChanged = 1;
    } else {
        for (int i = 0; i < g_stagingCount; i++) {
            FyneUIAElement* old = g_root->children[i];
            FyneUIAElement* neu = g_staging[i];
            if (old->controlType != neu->controlType ||
                wcscmp(old->name, neu->name) != 0) {
                structureChanged = 1;
                break;
            }
        }
    }

    if (structureChanged) {
        // Free old children
        for (int i = 0; i < g_root->childCount; i++) {
            free(g_root->children[i]->name);
            free(g_root->children[i]);
        }
        // Swap in staging
        FyneUIAElement** oldArr = g_root->children;
        int oldCap = g_root->childCapacity;
        g_root->children = g_staging;
        g_root->childCount = g_stagingCount;
        g_root->childCapacity = g_stagingCapacity;
        g_staging = oldArr;
        g_stagingCapacity = oldCap;
        g_stagingCount = 0;
        g_focusedIndex = -1;

        if (pfnUiaRaiseStructure) {
            int runtimeId[2] = { UiaAppendRuntimeId, g_root->uniqueId };
            pfnUiaRaiseStructure(&g_root->simple,
                UIAStructureChangeType_ChildrenInvalidated, runtimeId, 2);
        }
        PostMessageW(g_hwnd, WM_FYNE_RAISE_FOCUS, 0, 0);
    } else {
        // Update positions in-place
        for (int i = 0; i < g_stagingCount; i++) {
            FyneUIAElement* old = g_root->children[i];
            FyneUIAElement* neu = g_staging[i];
            old->x = neu->x;
            old->y = neu->y;
            old->width = neu->width;
            old->height = neu->height;
            free(neu->name);
            free(neu);
        }
        g_stagingCount = 0;
    }
}

void WinAccessibilityCleanup(void) {
    if (g_hwnd && g_origWndProc) {
        SetWindowLongPtrW(g_hwnd, GWLP_WNDPROC, (LONG_PTR)g_origWndProc);
        g_origWndProc = NULL;
    }
    if (g_root) {
        if (pfnUiaDisconnect) {
            pfnUiaDisconnect(&g_root->simple);
        }
        for (int i = 0; i < g_root->childCount; i++) {
            free(g_root->children[i]->name);
            free(g_root->children[i]);
        }
        free(g_root->children);
        free(g_root);
        g_root = NULL;
    }
    for (int i = 0; i < g_stagingCount; i++) {
        free(g_staging[i]->name);
        free(g_staging[i]);
    }
    free(g_staging);
    g_staging = NULL;
    g_stagingCount = 0;
    g_stagingCapacity = 0;
    g_hwnd = NULL;
}
