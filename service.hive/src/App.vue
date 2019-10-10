<template>
  <div id="app" class="container mx-auto bg-gray-200 p-3 pt-5 rounded-lg">
    <h1>Welcome to Your Vue.js App</h1>
    <button @click="getAssets" class="text-orange-500 hover:text-white hover:bg-orange-500 border border-orange-500 text-xs font-semibold rounded-full px-4 py-1 leading-normal">Refresh assets</button>
    <div class="flex flex-wrap justify-center">
      <HiveSwitch v-for="asset in assets" v-bind:key="asset.identifier" :asset="asset"></HiveSwitch>
    </div>
  </div>
</template>

<script>
  import HiveSwitch from './components/HiveSwitch.vue'

  export default {
    name: 'app',

    components: {
      HiveSwitch
    },

    mounted: function() {
      this.getAssets();
    },

    data: function() {
      return {
        assets: [],
      }
    },

    methods: {
      getAssets: function() {
        axios.get('http://' + host + ':8000/assets/')
          .then(response => this.assets = response.data);
      }
    }
  }
</script>

<style>

</style>
