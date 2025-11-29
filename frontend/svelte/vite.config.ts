import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [sveltekit()],
	server: {
		// Garante que o servidor escute fora do container
		host: '0.0.0.0', 
		port: 5175,
		
		allowedHosts: true
	}
});