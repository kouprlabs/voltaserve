// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

package processor

import (
	"fmt"

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
	err := infra.NewCommand().Exec("blender", "--background", "--python-expr", fmt.Sprintf(`
import bpy
import sys
from mathutils import Vector


def debug_print(message):
    print(message)


def get_combined_dimensions(objects):
    """
    Calculate the combined world-coordinate dimensions of the given objects.
    :param objects: List of Blender objects
    :return: Vector with dimensions (width, depth, height)
    """
    min_coord = Vector((float("inf"), float("inf"), float("inf")))
    max_coord = Vector((float("-inf"), float("-inf"), float("-inf")))

    for obj in objects:
        debug_print(f"Processing object: {obj.name}, type: {obj.type}")
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

    dimensions = max_coord - min_coord
    debug_print(f"Combined dimensions: {dimensions}")
    return dimensions


# Read command-line arguments
argv = sys.argv
argv = argv[argv.index("--") + 1 :]
input_file = argv[argv.index("--input") + 1]
output_file = argv[argv.index("--output") + 1]

# Clear existing objects
bpy.ops.object.select_all(action="SELECT")
bpy.ops.object.delete(use_global=False)

# Import GLB file
bpy.ops.import_scene.gltf(filepath=input_file)
debug_print("Imported GLB file")

# Collect all imported objects
imported_objects = bpy.context.selected_objects
if not imported_objects:
    raise RuntimeError("No objects imported")

# Move each imported object to center
for obj in imported_objects:
    obj.location = (0, 0, 0)
debug_print(f"Imported objects: {[obj.name for obj in imported_objects]}")

# Calculate combined dimensions for all mesh objects
dimensions = get_combined_dimensions(imported_objects)
max_dimension = max(dimensions)
debug_print(f"Max dimension: {max_dimension}")

# Determine a suitable distance for the camera based on the object's size
base_distance = 5.0
scaling_factor = 2.0  # Adjust this factor as needed
distance = base_distance + scaling_factor * max_dimension
debug_print(f"Calculated camera distance: {distance}")

# Add a camera
camera_data = bpy.data.cameras.new(name="Camera")
camera_object = bpy.data.objects.new("Camera", camera_data)
bpy.context.collection.objects.link(camera_object)
bpy.context.scene.camera = camera_object

# Set camera position and rotation
camera_object.location = (
    distance,
    -distance,
    distance * 0.65,
)  # Adjusted based on object size
camera_object.rotation_euler = (1.2, 0.0, 0.8)  # Maintain original rotation angles
debug_print(f"Camera location: {camera_object.location}")
debug_print(f"Camera rotation: {camera_object.rotation_euler}")

# Add a key light source
light_data = bpy.data.lights.new(name="KeyLight", type="POINT")
light_object = bpy.data.objects.new(name="KeyLight", object_data=light_data)
bpy.context.collection.objects.link(light_object)
light_object.location = (distance, -distance, distance)  # Adjusted based on object size
light_data.energy = 1500  # Increased intensity

# Add a fill light source (for better rendering)
fill_light_data = bpy.data.lights.new(name="FillLight", type="POINT")
fill_light_object = bpy.data.objects.new(name="FillLight", object_data=fill_light_data)
bpy.context.collection.objects.link(fill_light_object)
fill_light_object.location = (
    -distance,
    distance,
    distance,
)  # Adjusted based on object size
fill_light_data.energy = 1000  # Increased intensity
debug_print(f"Light locations: {light_object.location}, {fill_light_object.location}")

# Set white background
if bpy.context.scene.world.node_tree:
    world = bpy.context.scene.world
    nodes = world.node_tree.nodes
    background = nodes.get("Background", None)
    if background:
        background.inputs[0].default_value = (1, 1, 1, 1)  # White color (R, G, B, A)

# Set render settings
bpy.context.scene.render.engine = "CYCLES"
bpy.context.scene.render.filepath = output_file
bpy.context.scene.render.resolution_x = %d
bpy.context.scene.render.resolution_y = %d

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
debug_print("Render complete")`, width, height,
	), "--", "--input", inputPath, "--output", outputPath)
	if err != nil {
		return err
	}
	return nil
}
