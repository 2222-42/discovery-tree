# Production Build Guide

This document describes the production build configuration and deployment process for the React frontend application.

## Build Configuration

The production build is optimized with the following features:

### Code Splitting
- **Vendor chunk**: React and React-DOM libraries (11KB)
- **API chunk**: Axios and API utilities (35KB)
- **Main chunk**: Application code (208KB)
- **Dynamic chunks**: Configuration loaded on demand

### Bundle Optimization
- **Minification**: esbuild for fast, efficient minification
- **Tree shaking**: Automatic removal of unused code
- **CSS code splitting**: Separate CSS files for better caching
- **Asset organization**: JS, CSS, and images in separate directories
- **Source maps**: Generated for production debugging

### Performance Features
- **Modern target**: ES2020+ for better optimization
- **Compression**: Gzip compression analysis included
- **Manifest generation**: For deployment tools and CDN integration
- **Cache-friendly naming**: Hash-based file names for optimal caching

## Build Commands

### Standard Production Build
```bash
npm run build:prod
```
This command:
1. Cleans the previous build
2. Runs ESLint checks
3. Compiles TypeScript
4. Builds with production optimizations
5. Sets build timestamp

### Build Analysis
```bash
npm run build:analyze-only
```
Analyzes the current build output and provides:
- Bundle size breakdown
- Code splitting effectiveness
- Optimization recommendations
- Source map verification

### Full Build Test
```bash
npm run build:test
```
Comprehensive build validation that:
1. Performs a clean production build
2. Validates output structure
3. Analyzes bundle sizes
4. Tests preview server
5. Provides deployment readiness report

### Preview Production Build
```bash
npm run preview:prod
```
Builds and serves the production version locally for testing.

## Build Output Structure

```
dist/
├── index.html                    # Main HTML file
├── vite.svg                      # Static assets
├── .vite/
│   └── manifest.json            # Build manifest for deployment tools
└── assets/
    ├── css/
    │   └── index-[hash].css     # Compiled and minified CSS
    └── js/
        ├── vendor-[hash].js     # React/React-DOM bundle
        ├── api-[hash].js        # API utilities bundle
        ├── index-[hash].js      # Main application bundle
        ├── production-[hash].js # Production config (lazy loaded)
        └── *.js.map            # Source maps for debugging
```

## Environment Configuration

### Production Environment Variables
Set these in your deployment environment:

```bash
# API Configuration
VITE_API_BASE_URL=/api/v1

# Feature Flags
VITE_ENABLE_DEBUG_LOGGING=false
VITE_ENABLE_MOCK_API=false

# Build Information (optional)
VITE_APP_VERSION=1.0.0
VITE_BUILD_TIME=2024-01-01T00:00:00Z
```

### Production Config Features
The production configuration includes:
- API timeout and retry settings
- Performance optimizations
- Error handling configuration
- Feature flags for production
- Build metadata

## Deployment

### Static File Hosting
The build output is a standard static site that can be deployed to:
- **CDN**: CloudFront, CloudFlare, etc.
- **Static hosting**: Netlify, Vercel, GitHub Pages
- **Web servers**: Nginx, Apache
- **Cloud storage**: S3, Google Cloud Storage

### Server Configuration
For single-page application routing, configure your server to:
1. Serve `index.html` for all routes
2. Set appropriate cache headers for hashed assets
3. Enable gzip compression
4. Configure CORS if API is on different domain

### Example Nginx Configuration
```nginx
server {
    listen 80;
    server_name your-domain.com;
    root /path/to/dist;
    index index.html;

    # Cache static assets
    location /assets/ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }

    # SPA routing
    location / {
        try_files $uri $uri/ /index.html;
    }

    # Gzip compression
    gzip on;
    gzip_types text/css application/javascript application/json;
}
```

## Performance Metrics

### Current Build Metrics
- **Total bundle size**: ~280KB (optimal)
- **Main chunk**: 208KB (within recommended limits)
- **Vendor chunk**: 11KB (React optimized)
- **API chunk**: 35KB (utilities separated)
- **CSS**: 31KB (component styles)

### Optimization Recommendations
- ✅ Code splitting active
- ✅ Bundle sizes optimal
- ✅ Tree shaking enabled
- ✅ Source maps generated
- ✅ Asset organization implemented

## Monitoring and Debugging

### Source Maps
Source maps are generated for production debugging:
- Enable in browser dev tools
- Maps minified code back to original TypeScript
- Helps with error tracking and debugging

### Build Analysis
Regular build analysis helps maintain performance:
```bash
npm run build:analyze-only
```

### Performance Monitoring
Consider integrating:
- Web Vitals monitoring
- Bundle size tracking
- Performance budgets in CI/CD

## Troubleshooting

### Large Bundle Size
If bundles become too large:
1. Analyze with `npm run build:analyze-only`
2. Consider lazy loading more components
3. Review and remove unused dependencies
4. Split large components into separate chunks

### Build Failures
Common issues and solutions:
1. **TypeScript errors**: Fix type issues before building
2. **ESLint violations**: Run `npm run lint` to identify issues
3. **Memory issues**: Increase Node.js memory limit
4. **Asset loading**: Check asset paths and imports

### Runtime Issues
For production runtime issues:
1. Check browser console for errors
2. Verify API endpoints are accessible
3. Check network requests in dev tools
4. Use source maps for debugging minified code

## Security Considerations

### Environment Variables
- Never expose sensitive data in VITE_ variables
- Use server-side environment variables for secrets
- Validate all client-side configuration

### Content Security Policy
Consider implementing CSP headers:
```
Content-Security-Policy: default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline';
```

### Asset Integrity
- Use SRI (Subresource Integrity) for external resources
- Implement proper CORS policies
- Use HTTPS in production

## Continuous Integration

### Build Pipeline
Example CI/CD pipeline steps:
1. Install dependencies
2. Run linting and type checking
3. Run tests
4. Build production bundle
5. Analyze bundle size
6. Deploy to staging/production

### Quality Gates
- Bundle size limits
- Performance budgets
- Test coverage thresholds
- Security vulnerability scans