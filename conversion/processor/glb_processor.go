// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package processor

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/kouprlabs/voltaserve/conversion/config"
	"github.com/kouprlabs/voltaserve/conversion/infra"
)

type GLBProcessor struct {
	cmd       *infra.Command
	imageProc *ImageProcessor
	config    *config.Config
}

func NewGLBProcessor() *GLBProcessor {
	return &GLBProcessor{
		cmd:       infra.NewCommand(),
		imageProc: NewImageProcessor(),
		config:    config.GetConfig(),
	}
}

func (p *GLBProcessor) Thumbnail(inputPath string, width int, height int, outputPath string) error {
	err := infra.NewCommand().Exec("blender", "--background", "--python-expr", `
import bpy
import sys
from mathutils import Vector


def get_combined_dimensions_with_frames(objects):
    """
    Calculate the combined world-coordinate dimensions of the given objects,
    taking into account all animation frames.
    :param objects: List of Blender objects
    :return: Vector with dimensions (width, depth, height)
    """
    min_coord = Vector((float("inf"), float("inf"), float("inf")))
    max_coord = Vector((float("-inf"), float("-inf"), float("-inf")))

    # Get the current frame to restore it later
    current_frame = bpy.context.scene.frame_current

    # Iterate through all frames in the animation
    start_frame = bpy.context.scene.frame_start
    end_frame = bpy.context.scene.frame_end

    for frame in range(start_frame, end_frame + 1):
        bpy.context.scene.frame_set(frame)
        for obj in objects:
            if obj.type == "MESH":
                # Update the dependency graph
                depsgraph = bpy.context.evaluated_depsgraph_get()
                ob_eval = obj.evaluated_get(depsgraph)
                for vert in ob_eval.data.vertices:
                    world_coord = ob_eval.matrix_world @ vert.co
                    min_coord.x = min(min_coord.x, world_coord.x)
                    min_coord.y = min(min_coord.y, world_coord.y)
                    min_coord.z = min(min_coord.z, world_coord.z)
                    max_coord.x = max(max_coord.x, world_coord.x)
                    max_coord.y = max(max_coord.y, world_coord.y)
                    max_coord.z = max(max_coord.z, world_coord.z)

    # Restore the original frame
    bpy.context.scene.frame_set(current_frame)

    return max_coord - min_coord


def get_combined_dimensions(objects):
    """
    Calculate the combined world-coordinate dimensions of the given objects.
    :param objects: List of Blender objects
    :return: Vector with dimensions (width, depth, height)
    """
    min_coord = Vector((float("inf"), float("inf"), float("inf")))
    max_coord = Vector((float("-inf"), float("-inf"), float("-inf")))

    for obj in objects:
        if obj.type == "MESH":
            # Ensure transformations are applied
            bpy.context.view_layer.objects.active = obj
            bpy.ops.object.transform_apply(location=True, rotation=True, scale=True)

            # Update the dependency graph
            depsgraph = bpy.context.evaluated_depsgraph_get()
            ob_eval = obj.evaluated_get(depsgraph)

            for vert in ob_eval.data.vertices:
                world_coord = ob_eval.matrix_world @ vert.co
                min_coord.x = min(min_coord.x, world_coord.x)
                min_coord.y = min(min_coord.y, world_coord.y)
                min_coord.z = min(min_coord.z, world_coord.z)
                max_coord.x = max(max_coord.x, world_coord.x)
                max_coord.y = max(max_coord.y, world_coord.y)
                max_coord.z = max(max_coord.z, world_coord.z)

    return max_coord - min_coord


# Read command-line arguments
argv = sys.argv
argv = argv[argv.index("--") + 1 :]
input_file = argv[argv.index("--input") + 1]
output_file = argv[argv.index("--output") + 1]
has_animations = argv[argv.index("--animations") + 1].lower() == "true"
width = int(argv[argv.index("--width") + 1])
height = int(argv[argv.index("--height") + 1])

# Clear existing objects
bpy.ops.object.select_all(action="SELECT")
bpy.ops.object.delete(use_global=False)

# Import GLB file
bpy.ops.import_scene.gltf(filepath=input_file)

# Collect all imported objects
imported_objects = bpy.context.selected_objects
if not imported_objects:
    raise RuntimeError("No objects imported")

# Move each imported object to center
for obj in imported_objects:
    obj.location = (0, 0, 0)

# Calculate combined dimensions for all mesh objects
if has_animations:
    dimensions = get_combined_dimensions_with_frames(imported_objects)
else:
    dimensions = get_combined_dimensions(imported_objects)

# Determine a suitable distance for the camera based on the object's size
base_distance = 5.0
# Scaling factor can be adjusted if needed
scaling_factor = 2.0
# Calculate distance
distance = base_distance + scaling_factor * max(dimensions)

# Add a camera
camera_data = bpy.data.cameras.new(name="Camera")
camera_object = bpy.data.objects.new("Camera", camera_data)
bpy.context.collection.objects.link(camera_object)
bpy.context.scene.camera = camera_object

# Set camera position and rotation, adjusted based on object size
camera_object.location = (distance, -distance, distance * 0.65)
# Maintain original rotation angles
camera_object.rotation_euler = (1.2, 0.0, 0.8)

# Add a key light source
light_data = bpy.data.lights.new(name="KeyLight", type="POINT")
light_object = bpy.data.objects.new(name="KeyLight", object_data=light_data)
bpy.context.collection.objects.link(light_object)

# Adjusted based on object size
light_object.location = (distance, -distance, distance)
# Increased intensity
light_data.energy = 1500

# Add a fill light source (for better rendering)
fill_light_data = bpy.data.lights.new(name="FillLight", type="POINT")
fill_light_object = bpy.data.objects.new(name="FillLight", object_data=fill_light_data)
bpy.context.collection.objects.link(fill_light_object)

# Adjusted based on object size
fill_light_object.location = (-distance, distance, distance)

# Increased intensity
fill_light_data.energy = 1000

# Set white background
if bpy.context.scene.world.node_tree:
    world = bpy.context.scene.world
    nodes = world.node_tree.nodes
    background = nodes.get("Background", None)
    if background:
        # White color (R, G, B, A)
        background.inputs[0].default_value = (1, 1, 1, 1)

# Set render settings
bpy.context.scene.render.engine = "CYCLES"
bpy.context.scene.render.filepath = output_file
bpy.context.scene.render.resolution_x = width
bpy.context.scene.render.resolution_y = height

# Center and view all objects
bpy.ops.object.select_all(action="DESELECT")
for obj in imported_objects:
    obj.select_set(True)
bpy.context.view_layer.objects.active = imported_objects[0]
bpy.ops.view3d.camera_to_view_selected()

# Disable denoising because Blender on Debian/Ubuntu is built without OpenImageDenoiser
bpy.context.scene.cycles.use_denoising = False

# Render the scene
bpy.ops.render.render(write_still=True)
`,
		"--",
		"--input", inputPath,
		"--output", outputPath,
		"--animations", fmt.Sprintf("%t", p.hasAnimations(inputPath)),
		"--width", fmt.Sprintf("%d", width),
		"--height", fmt.Sprintf("%d", height),
	)
	if err != nil {
		return err
	}
	return nil
}

func (p *GLBProcessor) hasAnimations(filePath string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			infra.GetLogger().Error(err)
		}
	}(file)
	data, err := io.ReadAll(file)
	if err != nil {
		return false
	}
	// GLB Header is 12 bytes: magic (4 bytes) + version (4 bytes) + length (4 bytes)
	if len(data) < 12 || string(data[:4]) != "glTF" {
		return false
	}
	const headerSize = 12
	const chunkHeaderSize = 8
	offset := headerSize
	var jsonChunk []byte
	for offset < len(data) {
		if offset+chunkHeaderSize > len(data) {
			fmt.Println("Invalid chunk header in GLB file.")
			return false
		}
		chunkLength := binary.LittleEndian.Uint32(data[offset : offset+4])
		chunkType := string(data[offset+4 : offset+8])
		offset += chunkHeaderSize
		if offset+int(chunkLength) > len(data) {
			fmt.Println("Chunk length exceeds file size in GLB file.")
			return false
		}
		if chunkType == "JSON" {
			jsonChunk = data[offset : offset+int(chunkLength)]
			break
		}
		offset += int(chunkLength)
	}
	if jsonChunk == nil {
		return false
	}
	type GLTF struct {
		Animations []interface{} `json:"animations"`
	}
	var gltf GLTF
	if err := json.Unmarshal(jsonChunk, &gltf); err != nil {
		return false
	}
	return len(gltf.Animations) > 0
}
