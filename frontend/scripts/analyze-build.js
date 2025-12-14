#!/usr/bin/env node

/**
 * Build analysis script for production optimization validation
 */

import { existsSync, statSync, readFileSync, readdirSync } from 'fs';
import { join } from 'path';

const DIST_DIR = 'dist';
const MAX_CHUNK_SIZE = 500 * 1024; // 500KB
const MAX_TOTAL_SIZE = 2 * 1024 * 1024; // 2MB

console.log('📊 Analyzing production build...\n');

if (!existsSync(DIST_DIR)) {
  console.error('❌ Build directory not found! Run "npm run build" first.');
  process.exit(1);
}

// Check if manifest exists
const manifestPath = join(DIST_DIR, '.vite', 'manifest.json');
if (!existsSync(manifestPath)) {
  console.warn('⚠️  Manifest file not found');
} else {
  console.log('✅ Manifest file generated');
  
  // Analyze manifest
  const manifest = JSON.parse(readFileSync(manifestPath, 'utf8'));
  console.log('\n📋 Build Manifest Analysis:');
  
  Object.entries(manifest).forEach(([key, value]) => {
    if (value.isEntry) {
      console.log(`  🚪 Entry: ${key} -> ${value.file}`);
      if (value.imports) {
        console.log(`    📦 Imports: ${value.imports.length} chunks`);
      }
      if (value.dynamicImports) {
        console.log(`    🔄 Dynamic imports: ${value.dynamicImports.length} chunks`);
      }
      if (value.css) {
        console.log(`    🎨 CSS files: ${value.css.length}`);
      }
    } else if (value.isDynamicEntry) {
      console.log(`  🔄 Dynamic chunk: ${key} -> ${value.file}`);
    } else {
      console.log(`  📦 Chunk: ${value.name} -> ${value.file}`);
    }
  });
}

// Check main files exist
const indexPath = join(DIST_DIR, 'index.html');
if (!existsSync(indexPath)) {
  console.error('❌ index.html not found in build output!');
  process.exit(1);
}
console.log('\n✅ index.html generated');

// Analyze bundle sizes
console.log('\n📦 Bundle Size Analysis:');

const assetsDir = join(DIST_DIR, 'assets');
if (existsSync(assetsDir)) {
  let totalSize = 0;
  let hasLargeChunks = false;
  
  // Analyze JavaScript files
  const jsDir = join(assetsDir, 'js');
  if (existsSync(jsDir)) {
    console.log('\n  📄 JavaScript Files:');
    const jsFiles = readDirRecursive(jsDir).filter(f => f.endsWith('.js'));
    
    jsFiles.forEach(file => {
      const stats = statSync(file);
      const sizeKB = Math.round(stats.size / 1024);
      const fileName = file.split('/').pop();
      totalSize += stats.size;
      
      console.log(`    ${fileName}: ${sizeKB}KB`);
      
      if (stats.size > MAX_CHUNK_SIZE) {
        console.warn(`    ⚠️  Large chunk: ${sizeKB}KB (consider further code splitting)`);
        hasLargeChunks = true;
      }
    });
  }
  
  // Analyze CSS files
  const cssDir = join(assetsDir, 'css');
  if (existsSync(cssDir)) {
    console.log('\n  🎨 CSS Files:');
    const cssFiles = readDirRecursive(cssDir).filter(f => f.endsWith('.css'));
    
    cssFiles.forEach(file => {
      const stats = statSync(file);
      const sizeKB = Math.round(stats.size / 1024);
      const fileName = file.split('/').pop();
      totalSize += stats.size;
      
      console.log(`    ${fileName}: ${sizeKB}KB`);
    });
  }
  
  // Check for other assets
  const otherAssets = readDirRecursive(assetsDir).filter(f => 
    !f.endsWith('.js') && !f.endsWith('.css') && !f.endsWith('.map')
  );
  
  if (otherAssets.length > 0) {
    console.log('\n  🖼️  Other Assets:');
    otherAssets.forEach(file => {
      const stats = statSync(file);
      const sizeKB = Math.round(stats.size / 1024);
      const fileName = file.split('/').pop();
      totalSize += stats.size;
      
      console.log(`    ${fileName}: ${sizeKB}KB`);
    });
  }

  const totalSizeMB = (totalSize / (1024 * 1024)).toFixed(2);
  console.log(`\n📊 Total bundle size: ${totalSizeMB}MB`);

  // Provide optimization recommendations
  console.log('\n🎯 Optimization Analysis:');
  
  if (totalSize > MAX_TOTAL_SIZE) {
    console.warn(`  ⚠️  Large total bundle size: ${totalSizeMB}MB`);
    console.log('  💡 Consider:');
    console.log('     - Lazy loading more components');
    console.log('     - Removing unused dependencies');
    console.log('     - Using dynamic imports for large libraries');
  } else {
    console.log('  ✅ Total bundle size is optimal');
  }
  
  if (hasLargeChunks) {
    console.log('  💡 Consider splitting large chunks further');
  } else {
    console.log('  ✅ Individual chunk sizes are optimal');
  }
  
  // Check for code splitting effectiveness
  const jsFiles = readDirRecursive(assetsDir).filter(f => f.endsWith('.js'));
  if (jsFiles.length > 3) {
    console.log('  ✅ Code splitting is active');
  } else {
    console.log('  💡 Consider more aggressive code splitting');
  }
  
} else {
  console.warn('⚠️  Assets directory not found');
}

// Check for source maps
console.log('\n🗺️  Source Maps:');
const mapFiles = readDirRecursive(DIST_DIR).filter(f => f.endsWith('.map'));
if (mapFiles.length > 0) {
  console.log(`  ✅ ${mapFiles.length} source map files generated`);
} else {
  console.log('  ⚠️  No source maps found');
}

console.log('\n🎉 Build analysis complete!');

function readDirRecursive(dir) {
  const files = [];
  
  function traverse(currentDir) {
    const items = readdirSync(currentDir, { withFileTypes: true });
    
    for (const item of items) {
      const fullPath = join(currentDir, item.name);
      
      if (item.isDirectory()) {
        traverse(fullPath);
      } else {
        files.push(fullPath);
      }
    }
  }
  
  traverse(dir);
  return files;
}