<template><div>
	<web-header></web-header>
	<b-container id="app" fluid="lg">
		
		<br>
        <b-card header="Sunburst">
            <sunburst style="height: 500px; margin-bottom: -1em" :data="data">
				<breadcrumbTrail style="padding-top: 1em" slot="legend" slot-scope="{ nodes, colorGetter, width }" :current="nodes.mouseOver" :root="nodes.root" :colorGetter="colorGetter" :from="nodes.zoomed" :width="width" />
				
				<template slot-scope="{ on, actions }">
					<highlightOnHover v-bind="{ on, actions }"/>
					<zoomOnClick v-bind="{ on, actions }"/>
				</template>

            </sunburst>
        </b-card>
	</b-container>
</div></template>

<script>
import {
  breadcrumbTrail,
  highlightOnHover,
  nodeInfoDisplayer,
  sunburst,
  zoomOnClick
} from 'vue-d3-sunburst';
import "vue-d3-sunburst/dist/vue-d3-sunburst.css";

export default {
    components: {
    breadcrumbTrail,
    highlightOnHover,
    nodeInfoDisplayer,
    sunburst,
    zoomOnClick
  },
  data() {
    return {
      data:  {
      "name": "flare",
        "children": [
          {
            "name": "analytics",
            "children": [
              {
                "name": "cluster",
                "children": [
                  { "name": "AgglomerativeCluster", "size": 3938 },
                  { "name": "CommunityStructure", "size": 3812 },
                  { "name": "HierarchicalCluster", "size": 6714 },
                  { "name": "MergeEdge", "size": 743 }
                ]
              },
              {
                "name": "graph",
                "children": [
                  { "name": "BetweennessCentrality", "size": 3534 },
                  { "name": "LinkDistance", "size": 5731 },
                  { "name": "MaxFlowMinCut", "size": 7840 },
                  { "name": "ShortestPaths", "size": 5914 },
                  { "name": "SpanningTree", "size": 3416 }
                ]
              }
            ]
          }
        ]
      }
    }
  }
}
</script>