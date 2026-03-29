// Test script to check auth storage structure
// Run this in browser console after logging in

console.log('=== Checking Auth Storage ===');

// Check localStorage
const authData = localStorage.getItem('navhub-auth');
console.log('1. Raw localStorage data:', authData);

if (authData) {
  try {
    const parsed = JSON.parse(authData);
    console.log('2. Parsed structure:', JSON.stringify(parsed, null, 2));

    // Check different possible paths to token
    console.log('3. Checking token paths:');
    console.log('  - parsed.token:', parsed.token);
    console.log('  - parsed.state?.token:', parsed.state?.token);
    console.log('  - parsed.state?.state?.token:', parsed.state?.state?.token);

    // Check user data
    console.log('4. Checking user data:');
    console.log('  - parsed.user:', parsed.user);
    console.log('  - parsed.state?.user:', parsed.state?.user);

    // Check authentication status
    console.log('5. Authentication status:');
    console.log('  - parsed.isAuthenticated:', parsed.isAuthenticated);
    console.log('  - parsed.state?.isAuthenticated:', parsed.state?.isAuthenticated);
  } catch (e) {
    console.error('Error parsing auth data:', e);
  }
} else {
  console.log('No auth data found in localStorage');
}

// Check what the api interceptor would read
console.log('\n=== Simulating API Interceptor Logic ===');
const storedAuth = localStorage.getItem('navhub-auth');
if (storedAuth) {
  const auth = JSON.parse(storedAuth);
  console.log('api.ts reads: auth.state?.token =', auth.state?.token);

  if (auth.state?.token) {
    console.log('✅ Token would be set in Authorization header');
  } else {
    console.log('❌ Token would NOT be set - auth.state.token is undefined');
  }
}
