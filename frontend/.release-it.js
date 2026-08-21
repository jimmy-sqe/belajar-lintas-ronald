module.exports = {
  git: {
    commitMessage: 'chore: release ${version}',
    push: true,
    tagName: 'v${version}'
  },
  github: {
    release: true,
    web: true,
    releaseName: 'v${version}',
    releaseNotes(context) {
      return context.changelog;
    }
  },
  plugins: {
    '@release-it/conventional-changelog': {
      preset: {
        name: 'conventionalcommits',
        types: [
          { type: 'feat', section: 'Added' },
          { type: 'fix', section: 'Fixed' },
          { type: 'test', section: 'Test' },
          { type: 'perf', section: 'Improved' },
          { type: 'change', section: 'Changed' },
          { type: 'remove', section: 'Removed' }
        ]
      },
      infile: 'CHANGELOG.md',
      ignoreRecommendedBump: true
    }
  }
};
