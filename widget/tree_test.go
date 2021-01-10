package widget_test

import (
	"testing"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/stretchr/testify/assert"
)

var treeData = map[string][]string{
	"":    {"foo"},
	"foo": {"foobar"},
}

func TestTree_OpenClose(t *testing.T) {
	t.Run("Exists", func(t *testing.T) {
		tree := widget.NewTreeWithStrings(treeData)

		closed := make(chan string, 1)
		opened := make(chan string, 1)
		tree.OnBranchClosed = func(uid widget.TreeNodeID) {
			closed <- uid
		}
		tree.OnBranchOpened = func(uid widget.TreeNodeID) {
			opened <- uid
		}

		assert.False(t, tree.IsBranchOpen("foo"))

		tree.OpenBranch("foo")
		assert.True(t, tree.IsBranchOpen("foo"))

		select {
		case s := <-opened:
			assert.Equal(t, "foo", s)
		case <-time.After(1 * time.Second):
			assert.Fail(t, "Branch should have been opened")
		}

		tree.CloseBranch("foo")
		assert.False(t, tree.IsBranchOpen("foo"))

		select {
		case s := <-closed:
			assert.Equal(t, "foo", s)
		case <-time.After(1 * time.Second):
			assert.Fail(t, "Branch should have been closed")
		}

		tree.ToggleBranch("foo")
		assert.True(t, tree.IsBranchOpen("foo"))

		select {
		case s := <-opened:
			assert.Equal(t, "foo", s)
		case <-time.After(1 * time.Second):
			assert.Fail(t, "Branch should have been opened")
		}

		tree.ToggleBranch("foo")
		assert.False(t, tree.IsBranchOpen("foo"))

		select {
		case s := <-closed:
			assert.Equal(t, "foo", s)
		case <-time.After(1 * time.Second):
			assert.Fail(t, "Branch should have been closed")
		}
	})
	t.Run("Missing", func(t *testing.T) {
		tree := widget.NewTreeWithStrings(treeData)

		assert.False(t, tree.IsBranchOpen("foo"))

		tree.OpenBranch("bar")
		assert.False(t, tree.IsBranchOpen("foo"))

		tree.CloseBranch("bar")
		assert.False(t, tree.IsBranchOpen("foo"))

		tree.ToggleBranch("bar")
		assert.False(t, tree.IsBranchOpen("foo"))

		tree.ToggleBranch("bar")
		assert.False(t, tree.IsBranchOpen("foo"))
	})
}

func TestTree_OpenCloseAll(t *testing.T) {
	data := map[string][]string{
		"":     {"foo0", "foo1", "foo2"},
		"foo0": {"foobar0"},
		"foo1": {"foobar1"},
		"foo2": {"foobar2"},
	}
	tree := widget.NewTreeWithStrings(data)

	tree.OpenAllBranches()
	assert.True(t, tree.IsBranchOpen("foo0"))
	assert.True(t, tree.IsBranchOpen("foo1"))
	assert.True(t, tree.IsBranchOpen("foo2"))

	tree.CloseAllBranches()
	assert.False(t, tree.IsBranchOpen("foo0"))
	assert.False(t, tree.IsBranchOpen("foo1"))
	assert.False(t, tree.IsBranchOpen("foo2"))
}

func TestTree_Layout(t *testing.T) {
	test.NewApp()

	for name, tt := range map[string]struct {
		items    map[string][]string
		selected widget.TreeNodeID
		opened   []widget.TreeNodeID
		want     string
	}{
		"single_leaf": {
			items: map[string][]string{
				"": {"11111"},
			},
			want: `
				<canvas padded size="200x300">
					<content>
						<widget pos="4,4" size="192x292" type="*widget.Tree">
							<widget size="192x292" type="*widget.Scroll">
								<widget size="192x292" type="*widget.treeContent">
									<widget size="192x37" type="*widget.leaf">
										<widget pos="28,4" size="160x29" type="*widget.Label">
											<text pos="4,4" size="152x21">11111</text>
										</widget>
										<rectangle fillColor="background" size="4x37"/>
									</widget>
								</widget>
							</widget>
						</widget>
					</content>
				</canvas>
			`,
		},
		"single_leaf_selected": {
			items: map[string][]string{
				"": {"11111"},
			},
			selected: "11111",
			want: `
				<canvas padded size="200x300">
					<content>
						<widget pos="4,4" size="192x292" type="*widget.Tree">
							<widget size="192x292" type="*widget.Scroll">
								<widget size="192x292" type="*widget.treeContent">
									<widget size="192x37" type="*widget.leaf">
										<widget pos="28,4" size="160x29" type="*widget.Label">
											<text pos="4,4" size="152x21">11111</text>
										</widget>
										<rectangle fillColor="primary" size="4x37"/>
									</widget>
								</widget>
							</widget>
						</widget>
					</content>
				</canvas>
			`,
		},
		"single_branch": {
			items: map[string][]string{
				"":  {"A"},
				"A": {"11111"},
			},
			want: `
				<canvas padded size="200x300">
					<content>
						<widget pos="4,4" size="192x292" type="*widget.Tree">
							<widget size="192x292" type="*widget.Scroll">
								<widget size="192x292" type="*widget.treeContent">
									<widget size="192x37" type="*widget.branch">
										<widget pos="28,4" size="160x29" type="*widget.Label">
											<text pos="4,4" size="152x21">A</text>
										</widget>
										<widget pos="4,4" size="20x29" type="*widget.branchIcon">
											<image fillMode="contain" rsc="navigateNextIcon" size="20x29"/>
										</widget>
										<rectangle fillColor="background" size="4x37"/>
									</widget>
								</widget>
							</widget>
						</widget>
					</content>
				</canvas>
			`,
		},
		"single_branch_selected": {
			items: map[string][]string{
				"":  {"A"},
				"A": {"11111"},
			},
			selected: "A",
			want: `
				<canvas padded size="200x300">
					<content>
						<widget pos="4,4" size="192x292" type="*widget.Tree">
							<widget size="192x292" type="*widget.Scroll">
								<widget size="192x292" type="*widget.treeContent">
									<widget size="192x37" type="*widget.branch">
										<widget pos="28,4" size="160x29" type="*widget.Label">
											<text pos="4,4" size="152x21">A</text>
										</widget>
										<widget pos="4,4" size="20x29" type="*widget.branchIcon">
											<image fillMode="contain" rsc="navigateNextIcon" size="20x29"/>
										</widget>
										<rectangle fillColor="primary" size="4x37"/>
									</widget>
								</widget>
							</widget>
						</widget>
					</content>
				</canvas>
			`,
		},
		"single_branch_opened": {
			items: map[string][]string{
				"":  {"A"},
				"A": {"11111"},
			},
			opened: []string{"A"},
			want: `
				<canvas padded size="200x300">
					<content>
						<widget pos="4,4" size="192x292" type="*widget.Tree">
							<widget size="192x292" type="*widget.Scroll">
								<widget size="192x292" type="*widget.treeContent">
									<widget size="192x37" type="*widget.branch">
										<widget pos="28,4" size="160x29" type="*widget.Label">
											<text pos="4,4" size="152x21">A</text>
										</widget>
										<widget pos="4,4" size="20x29" type="*widget.branchIcon">
											<image fillMode="contain" rsc="moveDownIcon" size="20x29"/>
										</widget>
										<rectangle fillColor="background" size="4x37"/>
									</widget>
									<widget pos="4,37" size="184x1" type="*widget.Separator">
										<rectangle fillColor="disabled" size="184x1"/>
									</widget>
									<widget pos="0,38" size="192x37" type="*widget.leaf">
										<widget pos="52,4" size="136x29" type="*widget.Label">
											<text pos="4,4" size="128x21">11111</text>
										</widget>
										<rectangle fillColor="background" size="4x37"/>
									</widget>
								</widget>
							</widget>
						</widget>
					</content>
				</canvas>
			`,
		},
		"single_branch_opened_selected": {
			items: map[string][]string{
				"":  {"A"},
				"A": {"11111"},
			},
			opened:   []string{"A"},
			selected: "A",
			want: `
				<canvas padded size="200x300">
					<content>
						<widget pos="4,4" size="192x292" type="*widget.Tree">
							<widget size="192x292" type="*widget.Scroll">
								<widget size="192x292" type="*widget.treeContent">
									<widget size="192x37" type="*widget.branch">
										<widget pos="28,4" size="160x29" type="*widget.Label">
											<text pos="4,4" size="152x21">A</text>
										</widget>
										<widget pos="4,4" size="20x29" type="*widget.branchIcon">
											<image fillMode="contain" rsc="moveDownIcon" size="20x29"/>
										</widget>
										<rectangle fillColor="primary" size="4x37"/>
									</widget>
									<widget pos="4,37" size="184x1" type="*widget.Separator">
										<rectangle fillColor="disabled" size="184x1"/>
									</widget>
									<widget pos="0,38" size="192x37" type="*widget.leaf">
										<widget pos="52,4" size="136x29" type="*widget.Label">
											<text pos="4,4" size="128x21">11111</text>
										</widget>
										<rectangle fillColor="background" size="4x37"/>
									</widget>
								</widget>
							</widget>
						</widget>
					</content>
				</canvas>
			`,
		},
		"single_branch_opened_leaf_selected": {
			items: map[string][]string{
				"":  {"A"},
				"A": {"11111"},
			},
			opened:   []string{"A"},
			selected: "11111",
			want: `
				<canvas padded size="200x300">
					<content>
						<widget pos="4,4" size="192x292" type="*widget.Tree">
							<widget size="192x292" type="*widget.Scroll">
								<widget size="192x292" type="*widget.treeContent">
									<widget size="192x37" type="*widget.branch">
										<widget pos="28,4" size="160x29" type="*widget.Label">
											<text pos="4,4" size="152x21">A</text>
										</widget>
										<widget pos="4,4" size="20x29" type="*widget.branchIcon">
											<image fillMode="contain" rsc="moveDownIcon" size="20x29"/>
										</widget>
										<rectangle fillColor="background" size="4x37"/>
									</widget>
									<widget pos="4,37" size="184x1" type="*widget.Separator">
										<rectangle fillColor="disabled" size="184x1"/>
									</widget>
									<widget pos="0,38" size="192x37" type="*widget.leaf">
										<widget pos="52,4" size="136x29" type="*widget.Label">
											<text pos="4,4" size="128x21">11111</text>
										</widget>
										<rectangle fillColor="primary" size="4x37"/>
									</widget>
								</widget>
							</widget>
						</widget>
					</content>
				</canvas>
			`,
		},
		"multiple": {
			items: map[string][]string{
				"":  {"A", "B", "44444444444444444444"},
				"A": {"11111"},
				"B": {"2222222222"},
			},
			want: `
				<canvas padded size="200x300">
					<content>
						<widget pos="4,4" size="192x292" type="*widget.Tree">
							<widget size="192x292" type="*widget.Scroll">
								<widget size="192x292" type="*widget.treeContent">
									<widget size="192x37" type="*widget.branch">
										<widget pos="28,4" size="160x29" type="*widget.Label">
											<text pos="4,4" size="152x21">A</text>
										</widget>
										<widget pos="4,4" size="20x29" type="*widget.branchIcon">
											<image fillMode="contain" rsc="navigateNextIcon" size="20x29"/>
										</widget>
										<rectangle fillColor="background" size="4x37"/>
									</widget>
									<widget pos="4,37" size="184x1" type="*widget.Separator">
										<rectangle fillColor="disabled" size="184x1"/>
									</widget>
									<widget pos="0,38" size="192x37" type="*widget.branch">
										<widget pos="28,4" size="160x29" type="*widget.Label">
											<text pos="4,4" size="152x21">B</text>
										</widget>
										<widget pos="4,4" size="20x29" type="*widget.branchIcon">
											<image fillMode="contain" rsc="navigateNextIcon" size="20x29"/>
										</widget>
										<rectangle fillColor="background" size="4x37"/>
									</widget>
									<widget pos="4,75" size="184x1" type="*widget.Separator">
										<rectangle fillColor="disabled" size="184x1"/>
									</widget>
									<widget pos="0,76" size="192x37" type="*widget.leaf">
										<widget pos="28,4" size="160x29" type="*widget.Label">
											<text pos="4,4" size="152x21">44444444444444444444</text>
										</widget>
										<rectangle fillColor="background" size="4x37"/>
									</widget>
								</widget>
							</widget>
						</widget>
					</content>
				</canvas>
			`,
		},
		"multiple_selected": {
			items: map[string][]string{
				"":  {"A", "B", "44444444444444444444"},
				"A": {"11111"},
				"B": {"2222222222"},
			},
			selected: "44444444444444444444",
			want: `
				<canvas padded size="200x300">
					<content>
						<widget pos="4,4" size="192x292" type="*widget.Tree">
							<widget size="192x292" type="*widget.Scroll">
								<widget size="192x292" type="*widget.treeContent">
									<widget size="192x37" type="*widget.branch">
										<widget pos="28,4" size="160x29" type="*widget.Label">
											<text pos="4,4" size="152x21">A</text>
										</widget>
										<widget pos="4,4" size="20x29" type="*widget.branchIcon">
											<image fillMode="contain" rsc="navigateNextIcon" size="20x29"/>
										</widget>
										<rectangle fillColor="background" size="4x37"/>
									</widget>
									<widget pos="4,37" size="184x1" type="*widget.Separator">
										<rectangle fillColor="disabled" size="184x1"/>
									</widget>
									<widget pos="0,38" size="192x37" type="*widget.branch">
										<widget pos="28,4" size="160x29" type="*widget.Label">
											<text pos="4,4" size="152x21">B</text>
										</widget>
										<widget pos="4,4" size="20x29" type="*widget.branchIcon">
											<image fillMode="contain" rsc="navigateNextIcon" size="20x29"/>
										</widget>
										<rectangle fillColor="background" size="4x37"/>
									</widget>
									<widget pos="4,75" size="184x1" type="*widget.Separator">
										<rectangle fillColor="disabled" size="184x1"/>
									</widget>
									<widget pos="0,76" size="192x37" type="*widget.leaf">
										<widget pos="28,4" size="160x29" type="*widget.Label">
											<text pos="4,4" size="152x21">44444444444444444444</text>
										</widget>
										<rectangle fillColor="primary" size="4x37"/>
									</widget>
								</widget>
							</widget>
						</widget>
					</content>
				</canvas>
			`,
		},
		"multiple_leaf": {
			items: map[string][]string{
				"": {"11111", "2222222222", "333333333333333", "44444444444444444444"},
			},
			want: `
				<canvas padded size="200x300">
					<content>
						<widget pos="4,4" size="192x292" type="*widget.Tree">
							<widget size="192x292" type="*widget.Scroll">
								<widget size="192x292" type="*widget.treeContent">
									<widget size="192x37" type="*widget.leaf">
										<widget pos="28,4" size="160x29" type="*widget.Label">
											<text pos="4,4" size="152x21">11111</text>
										</widget>
										<rectangle fillColor="background" size="4x37"/>
									</widget>
									<widget pos="4,37" size="184x1" type="*widget.Separator">
										<rectangle fillColor="disabled" size="184x1"/>
									</widget>
									<widget pos="0,38" size="192x37" type="*widget.leaf">
										<widget pos="28,4" size="160x29" type="*widget.Label">
											<text pos="4,4" size="152x21">2222222222</text>
										</widget>
										<rectangle fillColor="background" size="4x37"/>
									</widget>
									<widget pos="4,75" size="184x1" type="*widget.Separator">
										<rectangle fillColor="disabled" size="184x1"/>
									</widget>
									<widget pos="0,76" size="192x37" type="*widget.leaf">
										<widget pos="28,4" size="160x29" type="*widget.Label">
											<text pos="4,4" size="152x21">333333333333333</text>
										</widget>
										<rectangle fillColor="background" size="4x37"/>
									</widget>
									<widget pos="4,113" size="184x1" type="*widget.Separator">
										<rectangle fillColor="disabled" size="184x1"/>
									</widget>
									<widget pos="0,114" size="192x37" type="*widget.leaf">
										<widget pos="28,4" size="160x29" type="*widget.Label">
											<text pos="4,4" size="152x21">44444444444444444444</text>
										</widget>
										<rectangle fillColor="background" size="4x37"/>
									</widget>
								</widget>
							</widget>
						</widget>
					</content>
				</canvas>
			`,
		},
		"multiple_leaf_selected": {
			items: map[string][]string{
				"": {"11111", "2222222222", "333333333333333", "44444444444444444444"},
			},
			selected: "2222222222",
			want: `
				<canvas padded size="200x300">
					<content>
						<widget pos="4,4" size="192x292" type="*widget.Tree">
							<widget size="192x292" type="*widget.Scroll">
								<widget size="192x292" type="*widget.treeContent">
									<widget size="192x37" type="*widget.leaf">
										<widget pos="28,4" size="160x29" type="*widget.Label">
											<text pos="4,4" size="152x21">11111</text>
										</widget>
										<rectangle fillColor="background" size="4x37"/>
									</widget>
									<widget pos="4,37" size="184x1" type="*widget.Separator">
										<rectangle fillColor="disabled" size="184x1"/>
									</widget>
									<widget pos="0,38" size="192x37" type="*widget.leaf">
										<widget pos="28,4" size="160x29" type="*widget.Label">
											<text pos="4,4" size="152x21">2222222222</text>
										</widget>
										<rectangle fillColor="primary" size="4x37"/>
									</widget>
									<widget pos="4,75" size="184x1" type="*widget.Separator">
										<rectangle fillColor="disabled" size="184x1"/>
									</widget>
									<widget pos="0,76" size="192x37" type="*widget.leaf">
										<widget pos="28,4" size="160x29" type="*widget.Label">
											<text pos="4,4" size="152x21">333333333333333</text>
										</widget>
										<rectangle fillColor="background" size="4x37"/>
									</widget>
									<widget pos="4,113" size="184x1" type="*widget.Separator">
										<rectangle fillColor="disabled" size="184x1"/>
									</widget>
									<widget pos="0,114" size="192x37" type="*widget.leaf">
										<widget pos="28,4" size="160x29" type="*widget.Label">
											<text pos="4,4" size="152x21">44444444444444444444</text>
										</widget>
										<rectangle fillColor="background" size="4x37"/>
									</widget>
								</widget>
							</widget>
						</widget>
					</content>
				</canvas>
			`,
		},
		"multiple_branch": {
			items: map[string][]string{
				"":  {"A", "B"},
				"A": {"11111"},
				"B": {"2222222222", "C"},
				"C": {"333333333333333"},
			},
			want: `
				<canvas padded size="200x300">
					<content>
						<widget pos="4,4" size="192x292" type="*widget.Tree">
							<widget size="192x292" type="*widget.Scroll">
								<widget size="192x292" type="*widget.treeContent">
									<widget size="192x37" type="*widget.branch">
										<widget pos="28,4" size="160x29" type="*widget.Label">
											<text pos="4,4" size="152x21">A</text>
										</widget>
										<widget pos="4,4" size="20x29" type="*widget.branchIcon">
											<image fillMode="contain" rsc="navigateNextIcon" size="20x29"/>
										</widget>
										<rectangle fillColor="background" size="4x37"/>
									</widget>
									<widget pos="4,37" size="184x1" type="*widget.Separator">
										<rectangle fillColor="disabled" size="184x1"/>
									</widget>
									<widget pos="0,38" size="192x37" type="*widget.branch">
										<widget pos="28,4" size="160x29" type="*widget.Label">
											<text pos="4,4" size="152x21">B</text>
										</widget>
										<widget pos="4,4" size="20x29" type="*widget.branchIcon">
											<image fillMode="contain" rsc="navigateNextIcon" size="20x29"/>
										</widget>
										<rectangle fillColor="background" size="4x37"/>
									</widget>
								</widget>
							</widget>
						</widget>
					</content>
				</canvas>
			`,
		},
		"multiple_branch_selected": {
			items: map[string][]string{
				"":  {"A", "B"},
				"A": {"11111"},
				"B": {"2222222222", "C"},
				"C": {"333333333333333"},
			},
			selected: "B",
			want: `
				<canvas padded size="200x300">
					<content>
						<widget pos="4,4" size="192x292" type="*widget.Tree">
							<widget size="192x292" type="*widget.Scroll">
								<widget size="192x292" type="*widget.treeContent">
									<widget size="192x37" type="*widget.branch">
										<widget pos="28,4" size="160x29" type="*widget.Label">
											<text pos="4,4" size="152x21">A</text>
										</widget>
										<widget pos="4,4" size="20x29" type="*widget.branchIcon">
											<image fillMode="contain" rsc="navigateNextIcon" size="20x29"/>
										</widget>
										<rectangle fillColor="background" size="4x37"/>
									</widget>
									<widget pos="4,37" size="184x1" type="*widget.Separator">
										<rectangle fillColor="disabled" size="184x1"/>
									</widget>
									<widget pos="0,38" size="192x37" type="*widget.branch">
										<widget pos="28,4" size="160x29" type="*widget.Label">
											<text pos="4,4" size="152x21">B</text>
										</widget>
										<widget pos="4,4" size="20x29" type="*widget.branchIcon">
											<image fillMode="contain" rsc="navigateNextIcon" size="20x29"/>
										</widget>
										<rectangle fillColor="primary" size="4x37"/>
									</widget>
								</widget>
							</widget>
						</widget>
					</content>
				</canvas>
			`,
		},
		"multiple_branch_opened": {
			items: map[string][]string{
				"":  {"A", "B"},
				"A": {"11111"},
				"B": {"2222222222", "C"},
				"C": {"333333333333333"},
			},
			opened: []string{"A", "B", "C"},
			want: `
				<canvas padded size="200x300">
					<content>
						<widget pos="4,4" size="192x292" type="*widget.Tree">
							<widget size="192x292" type="*widget.Scroll">
								<widget size="206x292" type="*widget.treeContent">
									<widget size="206x37" type="*widget.branch">
										<widget pos="28,4" size="174x29" type="*widget.Label">
											<text pos="4,4" size="166x21">A</text>
										</widget>
										<widget pos="4,4" size="20x29" type="*widget.branchIcon">
											<image fillMode="contain" rsc="moveDownIcon" size="20x29"/>
										</widget>
										<rectangle fillColor="background" size="4x37"/>
									</widget>
									<widget pos="4,37" size="198x1" type="*widget.Separator">
										<rectangle fillColor="disabled" size="198x1"/>
									</widget>
									<widget pos="0,38" size="206x37" type="*widget.leaf">
										<widget pos="52,4" size="150x29" type="*widget.Label">
											<text pos="4,4" size="142x21">11111</text>
										</widget>
										<rectangle fillColor="background" size="4x37"/>
									</widget>
									<widget pos="4,75" size="198x1" type="*widget.Separator">
										<rectangle fillColor="disabled" size="198x1"/>
									</widget>
									<widget pos="0,76" size="206x37" type="*widget.branch">
										<widget pos="28,4" size="174x29" type="*widget.Label">
											<text pos="4,4" size="166x21">B</text>
										</widget>
										<widget pos="4,4" size="20x29" type="*widget.branchIcon">
											<image fillMode="contain" rsc="moveDownIcon" size="20x29"/>
										</widget>
										<rectangle fillColor="background" size="4x37"/>
									</widget>
									<widget pos="4,113" size="198x1" type="*widget.Separator">
										<rectangle fillColor="disabled" size="198x1"/>
									</widget>
									<widget pos="0,114" size="206x37" type="*widget.leaf">
										<widget pos="52,4" size="150x29" type="*widget.Label">
											<text pos="4,4" size="142x21">2222222222</text>
										</widget>
										<rectangle fillColor="background" size="4x37"/>
									</widget>
									<widget pos="4,151" size="198x1" type="*widget.Separator">
										<rectangle fillColor="disabled" size="198x1"/>
									</widget>
									<widget pos="0,152" size="206x37" type="*widget.branch">
										<widget pos="52,4" size="150x29" type="*widget.Label">
											<text pos="4,4" size="142x21">C</text>
										</widget>
										<widget pos="28,4" size="20x29" type="*widget.branchIcon">
											<image fillMode="contain" rsc="moveDownIcon" size="20x29"/>
										</widget>
										<rectangle fillColor="background" size="4x37"/>
									</widget>
									<widget pos="4,189" size="198x1" type="*widget.Separator">
										<rectangle fillColor="disabled" size="198x1"/>
									</widget>
									<widget pos="0,190" size="206x37" type="*widget.leaf">
										<widget pos="76,4" size="126x29" type="*widget.Label">
											<text pos="4,4" size="118x21">333333333333333</text>
										</widget>
										<rectangle fillColor="background" size="4x37"/>
									</widget>
								</widget>
								<widget pos="0,286" size="192x6" type="*widget.scrollBarArea">
									<widget backgroundColor="scrollbar" pos="0,3" size="178x3" type="*widget.scrollBar">
									</widget>
								</widget>
								<widget pos="192,0" size="0x292" type="*widget.Shadow">
									<linearGradient angle="270" endColor="shadow" pos="-8,0" size="8x292"/>
								</widget>
							</widget>
						</widget>
					</content>
				</canvas>
			`,
		},
		"multiple_branch_opened_selected": {
			items: map[string][]string{
				"":  {"A", "B"},
				"A": {"11111"},
				"B": {"2222222222", "C"},
				"C": {"333333333333333"},
			},
			opened:   []string{"A", "B", "C"},
			selected: "B",
			want: `
				<canvas padded size="200x300">
					<content>
						<widget pos="4,4" size="192x292" type="*widget.Tree">
							<widget size="192x292" type="*widget.Scroll">
								<widget size="206x292" type="*widget.treeContent">
									<widget size="206x37" type="*widget.branch">
										<widget pos="28,4" size="174x29" type="*widget.Label">
											<text pos="4,4" size="166x21">A</text>
										</widget>
										<widget pos="4,4" size="20x29" type="*widget.branchIcon">
											<image fillMode="contain" rsc="moveDownIcon" size="20x29"/>
										</widget>
										<rectangle fillColor="background" size="4x37"/>
									</widget>
									<widget pos="4,37" size="198x1" type="*widget.Separator">
										<rectangle fillColor="disabled" size="198x1"/>
									</widget>
									<widget pos="0,38" size="206x37" type="*widget.leaf">
										<widget pos="52,4" size="150x29" type="*widget.Label">
											<text pos="4,4" size="142x21">11111</text>
										</widget>
										<rectangle fillColor="background" size="4x37"/>
									</widget>
									<widget pos="4,75" size="198x1" type="*widget.Separator">
										<rectangle fillColor="disabled" size="198x1"/>
									</widget>
									<widget pos="0,76" size="206x37" type="*widget.branch">
										<widget pos="28,4" size="174x29" type="*widget.Label">
											<text pos="4,4" size="166x21">B</text>
										</widget>
										<widget pos="4,4" size="20x29" type="*widget.branchIcon">
											<image fillMode="contain" rsc="moveDownIcon" size="20x29"/>
										</widget>
										<rectangle fillColor="primary" size="4x37"/>
									</widget>
									<widget pos="4,113" size="198x1" type="*widget.Separator">
										<rectangle fillColor="disabled" size="198x1"/>
									</widget>
									<widget pos="0,114" size="206x37" type="*widget.leaf">
										<widget pos="52,4" size="150x29" type="*widget.Label">
											<text pos="4,4" size="142x21">2222222222</text>
										</widget>
										<rectangle fillColor="background" size="4x37"/>
									</widget>
									<widget pos="4,151" size="198x1" type="*widget.Separator">
										<rectangle fillColor="disabled" size="198x1"/>
									</widget>
									<widget pos="0,152" size="206x37" type="*widget.branch">
										<widget pos="52,4" size="150x29" type="*widget.Label">
											<text pos="4,4" size="142x21">C</text>
										</widget>
										<widget pos="28,4" size="20x29" type="*widget.branchIcon">
											<image fillMode="contain" rsc="moveDownIcon" size="20x29"/>
										</widget>
										<rectangle fillColor="background" size="4x37"/>
									</widget>
									<widget pos="4,189" size="198x1" type="*widget.Separator">
										<rectangle fillColor="disabled" size="198x1"/>
									</widget>
									<widget pos="0,190" size="206x37" type="*widget.leaf">
										<widget pos="76,4" size="126x29" type="*widget.Label">
											<text pos="4,4" size="118x21">333333333333333</text>
										</widget>
										<rectangle fillColor="background" size="4x37"/>
									</widget>
								</widget>
								<widget pos="0,286" size="192x6" type="*widget.scrollBarArea">
									<widget backgroundColor="scrollbar" pos="0,3" size="178x3" type="*widget.scrollBar">
									</widget>
								</widget>
								<widget pos="192,0" size="0x292" type="*widget.Shadow">
									<linearGradient angle="270" endColor="shadow" pos="-8,0" size="8x292"/>
								</widget>
							</widget>
						</widget>
					</content>
				</canvas>
			`,
		},
		"multiple_branch_opened_leaf_selected": {
			items: map[string][]string{
				"":  {"A", "B"},
				"A": {"11111"},
				"B": {"2222222222", "C"},
				"C": {"333333333333333"},
			},
			opened:   []string{"A", "B", "C"},
			selected: "2222222222",
			want: `
				<canvas padded size="200x300">
					<content>
						<widget pos="4,4" size="192x292" type="*widget.Tree">
							<widget size="192x292" type="*widget.Scroll">
								<widget size="206x292" type="*widget.treeContent">
									<widget size="206x37" type="*widget.branch">
										<widget pos="28,4" size="174x29" type="*widget.Label">
											<text pos="4,4" size="166x21">A</text>
										</widget>
										<widget pos="4,4" size="20x29" type="*widget.branchIcon">
											<image fillMode="contain" rsc="moveDownIcon" size="20x29"/>
										</widget>
										<rectangle fillColor="background" size="4x37"/>
									</widget>
									<widget pos="4,37" size="198x1" type="*widget.Separator">
										<rectangle fillColor="disabled" size="198x1"/>
									</widget>
									<widget pos="0,38" size="206x37" type="*widget.leaf">
										<widget pos="52,4" size="150x29" type="*widget.Label">
											<text pos="4,4" size="142x21">11111</text>
										</widget>
										<rectangle fillColor="background" size="4x37"/>
									</widget>
									<widget pos="4,75" size="198x1" type="*widget.Separator">
										<rectangle fillColor="disabled" size="198x1"/>
									</widget>
									<widget pos="0,76" size="206x37" type="*widget.branch">
										<widget pos="28,4" size="174x29" type="*widget.Label">
											<text pos="4,4" size="166x21">B</text>
										</widget>
										<widget pos="4,4" size="20x29" type="*widget.branchIcon">
											<image fillMode="contain" rsc="moveDownIcon" size="20x29"/>
										</widget>
										<rectangle fillColor="background" size="4x37"/>
									</widget>
									<widget pos="4,113" size="198x1" type="*widget.Separator">
										<rectangle fillColor="disabled" size="198x1"/>
									</widget>
									<widget pos="0,114" size="206x37" type="*widget.leaf">
										<widget pos="52,4" size="150x29" type="*widget.Label">
											<text pos="4,4" size="142x21">2222222222</text>
										</widget>
										<rectangle fillColor="primary" size="4x37"/>
									</widget>
									<widget pos="4,151" size="198x1" type="*widget.Separator">
										<rectangle fillColor="disabled" size="198x1"/>
									</widget>
									<widget pos="0,152" size="206x37" type="*widget.branch">
										<widget pos="52,4" size="150x29" type="*widget.Label">
											<text pos="4,4" size="142x21">C</text>
										</widget>
										<widget pos="28,4" size="20x29" type="*widget.branchIcon">
											<image fillMode="contain" rsc="moveDownIcon" size="20x29"/>
										</widget>
										<rectangle fillColor="background" size="4x37"/>
									</widget>
									<widget pos="4,189" size="198x1" type="*widget.Separator">
										<rectangle fillColor="disabled" size="198x1"/>
									</widget>
									<widget pos="0,190" size="206x37" type="*widget.leaf">
										<widget pos="76,4" size="126x29" type="*widget.Label">
											<text pos="4,4" size="118x21">333333333333333</text>
										</widget>
										<rectangle fillColor="background" size="4x37"/>
									</widget>
								</widget>
								<widget pos="0,286" size="192x6" type="*widget.scrollBarArea">
									<widget backgroundColor="scrollbar" pos="0,3" size="178x3" type="*widget.scrollBar">
									</widget>
								</widget>
								<widget pos="192,0" size="0x292" type="*widget.Shadow">
									<linearGradient angle="270" endColor="shadow" pos="-8,0" size="8x292"/>
								</widget>
							</widget>
						</widget>
					</content>
				</canvas>
			`,
		},
	} {
		t.Run(name, func(t *testing.T) {
			tree := widget.NewTreeWithStrings(tt.items)
			for _, o := range tt.opened {
				tree.OpenBranch(o)
			}
			tree.Select(tt.selected)

			window := test.NewWindow(tree)
			defer window.Close()
			window.Resize(fyne.NewSize(200, 300))

			tree.Refresh() // Force layout

			test.AssertRendersToMarkup(t, tt.want, window.Canvas())
		})
	}
}

func TestTree_ChangeTheme(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	tree := widget.NewTreeWithStrings(treeData)
	tree.OpenBranch("foo")

	window := test.NewWindow(tree)
	defer window.Close()
	window.Resize(fyne.NewSize(220, 220))

	tree.Refresh() // Force layout

	test.AssertImageMatches(t, "tree/theme_initial.png", window.Canvas().Capture())

	test.WithTestTheme(t, func() {
		tree.Refresh()
		time.Sleep(100 * time.Millisecond)
		test.AssertImageMatches(t, "tree/theme_changed.png", window.Canvas().Capture())
	})
}

func TestTree_Move(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	tree := widget.NewTreeWithStrings(treeData)
	tree.OpenBranch("foo")

	window := test.NewWindow(tree)
	defer window.Close()
	window.Resize(fyne.NewSize(220, 220))

	tree.Refresh() // Force layout

	test.AssertRendersToMarkup(t, `
  	<canvas padded size="220x220">
  		<content>
  			<widget pos="4,4" size="212x212" type="*widget.Tree">
  				<widget size="212x212" type="*widget.Scroll">
  					<widget size="212x212" type="*widget.treeContent">
  						<widget size="212x37" type="*widget.branch">
  							<widget pos="28,4" size="180x29" type="*widget.Label">
  								<text pos="4,4" size="172x21">foo</text>
  							</widget>
  							<widget pos="4,4" size="20x29" type="*widget.branchIcon">
  								<image fillMode="contain" rsc="moveDownIcon" size="20x29"/>
  							</widget>
  							<rectangle fillColor="background" size="4x37"/>
  						</widget>
  						<widget pos="4,37" size="204x1" type="*widget.Separator">
  							<rectangle fillColor="disabled" size="204x1"/>
  						</widget>
  						<widget pos="0,38" size="212x37" type="*widget.leaf">
  							<widget pos="52,4" size="156x29" type="*widget.Label">
  								<text pos="4,4" size="148x21">foobar</text>
  							</widget>
  							<rectangle fillColor="background" size="4x37"/>
  						</widget>
  					</widget>
  				</widget>
  			</widget>
  		</content>
  	</canvas>
	`, window.Canvas())

	tree.Move(fyne.NewPos(20, 20))
	test.AssertRendersToMarkup(t, `
		<canvas padded size="220x220">
			<content>
				<widget pos="20,20" size="212x212" type="*widget.Tree">
					<widget size="212x212" type="*widget.Scroll">
						<widget size="212x212" type="*widget.treeContent">
							<widget size="212x37" type="*widget.branch">
								<widget pos="28,4" size="180x29" type="*widget.Label">
									<text pos="4,4" size="172x21">foo</text>
								</widget>
								<widget pos="4,4" size="20x29" type="*widget.branchIcon">
									<image fillMode="contain" rsc="moveDownIcon" size="20x29"/>
								</widget>
								<rectangle fillColor="background" size="4x37"/>
							</widget>
							<widget pos="4,37" size="204x1" type="*widget.Separator">
								<rectangle fillColor="disabled" size="204x1"/>
							</widget>
							<widget pos="0,38" size="212x37" type="*widget.leaf">
								<widget pos="52,4" size="156x29" type="*widget.Label">
									<text pos="4,4" size="148x21">foobar</text>
								</widget>
								<rectangle fillColor="background" size="4x37"/>
							</widget>
						</widget>
					</widget>
				</widget>
			</content>
		</canvas>
`, window.Canvas())
}

func TestTree_Refresh(t *testing.T) {
	app := test.NewApp()
	defer test.NewApp()
	app.Settings().SetTheme(theme.LightTheme())

	value := "Foo Leaf"
	tree := widget.NewTreeWithStrings(treeData)
	tree.UpdateNode = func(uid widget.TreeNodeID, branch bool, node fyne.CanvasObject) {
		if uid == "foobar" {
			node.(*widget.Label).SetText(value)
			assert.False(t, branch)
		} else {
			node.(*widget.Label).SetText(uid)
			assert.True(t, branch)
		}
	}
	tree.OpenBranch("foo")

	window := test.NewWindow(tree)
	defer window.Close()
	window.Resize(fyne.NewSize(220, 220))

	test.AssertImageMatches(t, "tree/refresh_initial.png", window.Canvas().Capture())

	value = "Replaced"
	tree.Refresh()

	test.AssertImageMatches(t, "tree/refresh_replaced.png", window.Canvas().Capture())
}
