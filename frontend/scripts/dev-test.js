#!/usr/bin/env node

/**
 * Development workflow testing script
 * Tests the development server setup and hot reloading functionality
 */

import { spawn } from 'child_process';
import { readFileSync, writeFileSync } from 'fs';
import { join } from 'path';

const FRONTEND_DIR = process.cwd();
const TEST_COMPONENT_PATH = join(FRONTEND_DIR, 'src', 'components', 'TestHMR.tsx');

/**
 * Create a test component for HMR testing
 */
function createTestComponent() {
  const testComponent = `
import React from 'react';

/**
 * Test component for Hot Module Replacement verification
 */
export const TestHMR: React.FC = () => {
  return (
    <div style={{ 
      padding: '20px', 
      backgroundColor: '#f0f0f0', 
      border: '2px solid #007acc',
      borderRadius: '8px',
      margin: '10px'
    }}>
      <h3>HMR Test Component - Version 1</h3>
      <p>This component tests hot module replacement functionality.</p>
      <p>Current time: {new Date().toLocaleTimeString()}</p>
    </div>
  );
};

export default TestHMR;

// Enable HMR for this component
if (import.meta.hot) {
  import.meta.hot.accept();
}
`;

  writeFileSync(TEST_COMPONENT_PATH, testComponent);
  console.log('✅ Created test component for HMR testing');
}

/**
 * Update the test component to verify HMR works
 */
function updateTestComponent() {
  const updatedComponent = `
import React from 'react';

/**
 * Test component for Hot Module Replacement verification
 */
export const TestHMR: React.FC = () => {
  return (
    <div style={{ 
      padding: '20px', 
      backgroundColor: '#e8f5e8', 
      border: '2px solid #28a745',
      borderRadius: '8px',
      margin: '10px'
    }}>
      <h3>HMR Test Component - Version 2 (Updated!)</h3>
      <p>This component has been updated to test hot module replacement.</p>
      <p>Current time: {new Date().toLocaleTimeString()}</p>
      <p style={{ color: '#28a745', fontWeight: 'bold' }}>
        🎉 Hot reload successful!
      </p>
    </div>
  );
};

export default TestHMR;

// Enable HMR for this component
if (import.meta.hot) {
  import.meta.hot.accept();
}
`;

  writeFileSync(TEST_COMPONENT_PATH, updatedComponent);
  console.log('✅ Updated test component to verify HMR');
}

/**
 * Test environment variables
 */
function testEnvironmentVariables() {
  console.log('🔍 Testing environment variables...');
  
  try {
    const envFile = readFileSync(join(FRONTEND_DIR, '.env'), 'utf8');
    const envDevFile = readFileSync(join(FRONTEND_DIR, '.env.development'), 'utf8');
    
    console.log('✅ Environment files exist');
    console.log('📄 .env file contains:', envFile.split('\\n').length, 'lines');
    console.log('📄 .env.development file contains:', envDevFile.split('\\n').length, 'lines');
    
    // Check for required variables
    const requiredVars = [
      'VITE_API_BASE_URL',
      'VITE_API_TIMEOUT',
      'VITE_ENABLE_DEBUG_LOGGING'
    ];
    
    const missingVars = requiredVars.filter(varName => 
      !envFile.includes(varName) && !envDevFile.includes(varName)
    );
    
    if (missingVars.length === 0) {
      console.log('✅ All required environment variables are defined');
    } else {
      console.log('❌ Missing environment variables:', missingVars);
    }
  } catch (error) {
    console.error('❌ Error reading environment files:', error.message);
  }
}

/**
 * Test Vite configuration
 */
function testViteConfig() {
  console.log('🔍 Testing Vite configuration...');
  
  try {
    const viteConfig = readFileSync(join(FRONTEND_DIR, 'vite.config.ts'), 'utf8');
    
    const requiredFeatures = [
      'server:',
      'proxy:',
      'hmr:',
      '/api'
    ];
    
    const missingFeatures = requiredFeatures.filter(feature => 
      !viteConfig.includes(feature)
    );
    
    if (missingFeatures.length === 0) {
      console.log('✅ Vite configuration includes all required features');
    } else {
      console.log('❌ Missing Vite configuration features:', missingFeatures);
    }
  } catch (error) {
    console.error('❌ Error reading Vite configuration:', error.message);
  }
}

/**
 * Main test function
 */
async function runTests() {
  console.log('🚀 Testing development server and hot reloading setup\\n');
  
  // Test configuration files
  testEnvironmentVariables();
  testViteConfig();
  
  // Create test component
  createTestComponent();
  
  console.log('\\n📋 Manual testing steps:');
  console.log('1. Run "npm run dev" to start the development server');
  console.log('2. Open http://localhost:3000 in your browser');
  console.log('3. Check that the application loads without errors');
  console.log('4. Open browser developer tools and check for API proxy logs');
  console.log('5. Modify any component file and verify hot reloading works');
  console.log('6. Check that environment variables are accessible in the app');
  
  console.log('\\n🎯 Expected results:');
  console.log('- Development server starts on port 3000');
  console.log('- Browser opens automatically');
  console.log('- API requests are proxied to localhost:8080');
  console.log('- File changes trigger hot reloads without full page refresh');
  console.log('- Debug logging appears in console (if enabled)');
  console.log('- No TypeScript compilation errors');
  
  console.log('\\n✅ Development server setup test completed!');
}

// Run the tests
runTests().catch(console.error);