<template>
  <div class="min-h-screen bg-gray-900 text-white font-sans">
    <nav class="bg-gray-800 p-4 shadow-lg">
      <div class="container mx-auto flex justify-between items-center">
        <h1 class="text-xl font-bold text-purple-500">Game Tracker</h1>
        <div class="space-x-4">
          <router-link to="/playing" class="hover:text-purple-400" active-class="text-purple-400 font-bold">Playing</router-link>
          <router-link to="/backlog" class="hover:text-purple-400" active-class="text-purple-400 font-bold">Backlog</router-link>
          <router-link to="/history" class="hover:text-purple-400" active-class="text-purple-400 font-bold">History</router-link>
          <router-link to="/calendar" class="hover:text-purple-400" active-class="text-purple-400 font-bold">Calendar</router-link>
        </div>
      </div>
    </nav>
    <main class="container mx-auto p-4">
      <router-view></router-view>
    </main>
  </div>
</template>
