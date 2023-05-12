<template>
  <div id="app">
    <div class="row">
      <div class="col col-4-lg">
        <fieldset>
          <legend>Legend</legend>
          <b>Instructions</b>
          <hr />
          <p>
            Click or hover on node to isolate that node's connections. Click on
            the light purple background to unselect.
          </p>
          <p>
            All resources that the node depends on are represented by a solid
            line. All resources that depend on the node are represented by a
            dashed line.
          </p>
          <hr />
          <b>Resource</b>
          <hr />
          <div class="node create">Resource - Create</div>
          <div class="node delete">Resource - Delete</div>
          <div class="node replace">Resource - Replace</div>
          <div class="node update">Resource - Update</div>
          <div class="node no-op">Resource - No Operation</div>
          <hr />
          <b>Other items</b>
          <hr />
          <div class="node variable">Variable</div>
          <div class="node output">Output</div>
          <div class="node data">Data</div>
          <div class="node module">Module</div>
          <div class="node locals">Local</div>
          <hr />
        </fieldset>
        <resource-detail :resourceID="resourceID" />
      </div>
      <div class="col col-8-lg">
        <graph ref="filegraph" :displayGraph="displayGraph" v-on:getNode="selectResource" />
        <explorer @selectResource="selectResource" />
      </div>
    </div>
  </div>
</template>

<script>
import ResourceDetail from "@/components/ResourceDetail.vue";
import Explorer from "@/components/Explorer.vue";

export default {
  name: "App",
  metaInfo: {
    title: "Rover | Terraform Visualization",
  },
  components: {
    Explorer,
    ResourceDetail,
  },
  data() {
    return {
      displayGraph: true,
      resourceID: "",
    };
  },
  methods: {
    saveGraph() {
      // this.displayGraph = displayGraph;
      this.$refs.filegraph.saveGraph();
    },
    selectResource(resourceID) {
      this.resourceID = resourceID;
    },
  },
};
</script>

<style scoped>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  /* text-align: center; */
  margin: 0 auto;
  margin-top: 60px;
  width: 90%;
}

.node {
  display: inline-block;
  margin: 0 1%;
  width: 48%;
  font-size: 0.9em;
}

.node:hover {
  transform: scale(1.02);
}

.resource-type {
  width: 20em;
  font-size: 2em;
  height: 100%;
}

.create {
  background-color: #28a745;
  color: white;
  font-weight: bold;
  border: 0;
}

.delete {
  /* background-color: #ffe9e9;
  border: 5px solid #e40707; */
  background-color: #e40707;
  color: white;
  font-weight: bold;
  border: 0;
}

.update {
  /* background-color: #e1f0ff;
  border: 5px solid #1d7ada; */
  background-color: #1d7ada;
  color: white;
  font-weight: bold;
  border: 0;
}

.replace {
  /* background-color: #fff7e0;
  border: 5px solid #ffc107; */
  background-color: #ffc107;
  color: black;
  font-weight: bold;
  border: 0;
}

.output {
  background-color: #fff7e0;
  border: 5px solid #ffc107;
  color: black;
  font-weight: bold;
}

.variable {
  background-color: #e1f0ff;
  border: 5px solid #1d7ada;
  color: black;
  font-weight: bold;
}

.data {
  background-color: #ffecec;
  border: 5px solid #dc477d;
  color: black;
  font-weight: bold;
}

.locals {
  background-color: black;
  color: white;
  font-weight: bold;
  border: 0;
}

.module {
  border: 5px solid #8450ba;
  color: #8450ba;
}
</style>
