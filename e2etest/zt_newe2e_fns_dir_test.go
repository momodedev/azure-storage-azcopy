package e2etest

import (
	"github.com/Azure/azure-storage-azcopy/v10/common"
	"time"
)

/*
FNSSuite exists to test oddities about virtual directory semantics on flat namespace blob.
*/
type FNSSuite struct{}

func init() {
	suiteManager.RegisterSuite(&FNSSuite{})
}

func (*FNSSuite) Scenario_CopyToOverlappableDirectoryMarker(a *ScenarioVariationManager) {
	tgtVerb := ResolveVariation(a, []AzCopyVerb{AzCopyVerbCopy, AzCopyVerbSync})

	// Target a fns account
	destRm := ObjectResourceMappingFlat{
		"foobar/": ResourceDefinitionObject{
			ObjectProperties: ObjectProperties{
				Metadata: common.Metadata{
					common.POSIXFolderMeta: pointerTo("true"),
				},
			},
			Body: NewZeroObjectContentContainer(0),
		},
	}

	if tgtVerb == AzCopyVerbSync {
		// Sync must have an existing destination, non-folder.
		destRm["foobar"] = ResourceDefinitionObject{
			Body: NewZeroObjectContentContainer(512),
		}
	}

	dest := CreateResource[ContainerResourceManager](a, GetRootResource(a, common.ELocation.Blob()),
		ResourceDefinitionContainer{
			Objects: destRm,
		},
	)

	if tgtVerb == AzCopyVerbSync && !a.Dryrun() {
		time.Sleep(time.Second * 5) // Ensure the source is newer
	}

	// Source must be newer than the destination
	source := CreateResource[ObjectResourceManager](a, GetRootResource(a, common.ELocation.Local()), ResourceDefinitionObject{
		Body: NewRandomObjectContentContainer(1024),
	})

	_, _ = RunAzCopy(a,
		AzCopyCommand{
			Verb: tgtVerb,
			Targets: []ResourceManager{
				source,
				dest.GetObject(a, "foobar", common.EEntityType.File()),
			},
			Flags: CopyFlags{
				AsSubdir: common.Iff(tgtVerb == AzCopyVerbCopy, pointerTo(false), nil),
			},
		},
	)

	ValidateResource(a, dest, ResourceDefinitionContainer{
		Objects: ObjectResourceMappingFlat{
			"foobar": ResourceDefinitionObject{
				ObjectShouldExist: pointerTo(true),
			},
			"foobar/": ResourceDefinitionObject{
				ObjectProperties: ObjectProperties{
					Metadata: common.Metadata{
						common.POSIXFolderMeta: pointerTo("true"),
					},
				},
				ObjectShouldExist: pointerTo(true),
			},
		},
	}, true)
}

// Scenario_IncludeRootDirectoryStub tests that the root directory (and sub directories) appropriately get their files picked up.
func (*FNSSuite) Scenario_IncludeRootDirectoryStub(a *ScenarioVariationManager) {
	dst := CreateResource[ContainerResourceManager](a, GetRootResource(a, common.ELocation.Blob()), ResourceDefinitionContainer{})
	src := CreateResource[ContainerResourceManager](a, GetRootResource(a, common.ELocation.Blob()), ResourceDefinitionContainer{
		Objects: ObjectResourceMappingFlat{
			"foobar":               ResourceDefinitionObject{Body: NewRandomObjectContentContainer(512), ObjectProperties: ObjectProperties{Metadata: common.Metadata{"dontcopyme": pointerTo("")}}}, // Object w/ same name as root dir
			"foobar/":              ResourceDefinitionObject{ObjectProperties: ObjectProperties{EntityType: common.EEntityType.Folder(), Metadata: common.Metadata{"asdf": pointerTo("qwerty")}}},    // Folder w/ same name as object, add special prop to ensure
			"foobar/foo":           ResourceDefinitionObject{Body: NewZeroObjectContentContainer(0)},
			"foobar/bar":           ResourceDefinitionObject{Body: NewZeroObjectContentContainer(0)},
			"foobar/baz":           ResourceDefinitionObject{Body: NewZeroObjectContentContainer(0)},
			"foobar/folder/":       ResourceDefinitionObject{ObjectProperties: ObjectProperties{EntityType: common.EEntityType.Folder()}},
			"foobar/folder/foobar": ResourceDefinitionObject{Body: NewZeroObjectContentContainer(0)},
		},
	})

	azcopyVerb := ResolveVariation(a, []AzCopyVerb{AzCopyVerbCopy, AzCopyVerbSync})
	RunAzCopy(a,
		AzCopyCommand{
			Verb: azcopyVerb,
			Targets: []ResourceManager{
				src.GetObject(a, "foobar/", common.EEntityType.Folder()),
				dst.GetObject(a, "foobar/", common.EEntityType.Folder()),
			},
			Flags: CopyFlags{
				CopySyncCommonFlags: CopySyncCommonFlags{
					Recursive:             pointerTo(true),
					IncludeDirectoryStubs: pointerTo(true),
				},
				AsSubdir: common.Iff(azcopyVerb == AzCopyVerbCopy, pointerTo(false), nil),
			},
		},
	)

	ValidateResource(a, dst, ResourceDefinitionContainer{
		Objects: ObjectResourceMappingFlat{
			"foobar":               ResourceDefinitionObject{ObjectShouldExist: pointerTo(false)}, // We shouldn't have captured foobar, but foobar/ should exist as a directory.
			"foobar/":              ResourceDefinitionObject{ObjectProperties: ObjectProperties{EntityType: common.EEntityType.Folder(), Metadata: common.Metadata{common.POSIXFolderMeta: pointerTo("true"), "asdf": pointerTo("qwerty")}}},
			"foobar/foo":           ResourceDefinitionObject{Body: NewZeroObjectContentContainer(0)},
			"foobar/bar":           ResourceDefinitionObject{Body: NewZeroObjectContentContainer(0)},
			"foobar/baz":           ResourceDefinitionObject{Body: NewZeroObjectContentContainer(0)},
			"foobar/folder/":       ResourceDefinitionObject{ObjectProperties: ObjectProperties{EntityType: common.EEntityType.Folder()}},
			"foobar/folder/foobar": ResourceDefinitionObject{Body: NewZeroObjectContentContainer(0)},
		},
	}, false)
}
