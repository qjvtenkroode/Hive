<template>
  <div>
    <form id="registry">
      <div>
        <label for="id">Identifier</label>
        <input id="id" v-model="id" type="text" name="id">
      </div>
      <div>
        <label for="name">Name</label>
        <input id="name" v-model="name" type="text" name="name">
      </div>
      <div>
        <label for="type">Type</label>
        <input id="type" v-model="type" type="text" name="type">
      </div>
      <div>
        <label for="controller">Controller</label>
        <input id="controller" v-model="controller" type="text" name="controller">
      </div>
      <div>
        <button @click="postAsset" type="button">Add Asset</button>
        <button @click="deleteAsset" type="button">Delete Asset</button>
      </div>
    </form>
  </div>
</template>

<script>
  export default {
    name: 'HiveRegistry',

    data: function() {
      return {
        id: null,
        name: null,
        type: null,
        controller: null,
      }
    },

    methods: {
      postAsset: function() {
        var asset = { identifier: this.id, name: this.name, type: this.type, controller: this.controller };
        console.log(asset);
        axios.post('http://' + host + ':8000/assets/', asset);
        this.clearFields()
      },
      
      deleteAsset: function() {
        var asset = { identifier: this.id };
        console.log(asset);
        axios.delete('http://' + host + ':8000/assets/' + asset.identifier, asset);
        this.clearFields()
      },

      clearFields: function() {
        this.id = null;
        this.name = null;
        this.type = null;
        this.controller = null;
      }
    }
  }
</script>

<style scoped>
</style>
