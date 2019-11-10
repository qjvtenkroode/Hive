<template>
  <div id="app" class="container mx-auto bg-gray-200 p-3 pt-5 rounded-lg">
    <ul>
      <li><router-link to="/">Switches</router-link></li>
      <li><router-link to="/rooms">Rooms</router-link></li>
      <li><router-link to="/registry">Registry</router-link></li>
    </ul>
    <router-view></router-view>
  </div>
</template>

<script>

</script>

<style>

</style>
