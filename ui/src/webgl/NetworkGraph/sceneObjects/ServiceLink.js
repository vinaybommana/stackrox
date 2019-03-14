import * as THREE from 'three';
import { MeshLine, MeshLineMaterial } from 'three.meshline';
import * as constants from 'constants/networkGraph';

const ServiceLink = (scene, canvas, data) => {
    const link = data;

    let line = null;
    let geometry = null;
    let meshLine = null;
    let material = null;

    function createLink() {
        // create a link mesh
        geometry = new THREE.Geometry();
        geometry.vertices[0] = new THREE.Vector3(0, 0, 0);
        geometry.vertices[1] = new THREE.Vector3(0, 0, 0);
        geometry.verticesNeedUpdate = true;
        /*
         * Using MeshLine instead of Line because due to limitations of the OpenGL Core Profile with
         * the WebGL renderer on most platforms linewidth will always be 1 regardless of the set value.
         */
        meshLine = new MeshLine();
        meshLine.setGeometry(geometry);
        const materialConfig = {
            useMap: false,
            color: new THREE.Color(constants.LINK_COLOR),
            opacity: 1,
            transparent: true,
            resolution: new THREE.Vector2(canvas.clientWidth, canvas.clientHeight),
            sizeAttenuation: true,
            lineWidth: constants.NODE_LINK_WIDTH
        };

        if (!link.isActive) {
            materialConfig.dashArray = constants.NODE_DASH_ARRAY;
            materialConfig.dashOffset = constants.NODE_DASH_OFFSET;
            materialConfig.dashRatio = constants.NODE_DASH_RATIO;
        }
        material = new MeshLineMaterial(materialConfig);

        line = new THREE.Mesh(meshLine.geometry, material);
        line.frustumCulled = false;
        line.mLine = meshLine;
        line.geo = geometry;
        line.userData = { link };

        scene.add(line);
    }

    function removeLink() {
        if (!line) return;
        scene.remove(line);
        line.geometry.dispose();
        meshLine.geometry.dispose();
        line.material.dispose();
        geometry.dispose();
        material.dispose();
        line = null;
    }

    function getSource() {
        return link.source;
    }

    function getTarget() {
        return link.target;
    }

    function isLinkInScene() {
        return !!line;
    }

    function update() {
        const { source, target } = link;
        if (line) {
            line.geo.vertices[0].x = source.x;
            line.geo.vertices[0].y = source.y;
            line.geo.vertices[1].x = target.x;
            line.geo.vertices[1].y = target.y;
            line.mLine.setGeometry(line.geo);
        }
    }

    function getType() {
        return constants.NETWORK_GRAPH_TYPES.LINK;
    }

    function cleanUp() {
        removeLink();
    }

    return {
        update,
        removeLink,
        createLink,
        getSource,
        getTarget,
        isLinkInScene,
        getType,
        cleanUp
    };
};

export default ServiceLink;
