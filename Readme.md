gh api \
  --method POST \
  -H "Accept: application/vnd.github+json" \
  -H "X-GitHub-Api-Version: 2022-11-28" \
  /repos/godarkthrt/cmdHelperTools/releases \
  -f tag_name='v1.0.0' \
 -f target_commitish='main' \
 -f name='v1.0.0' \
 -f body='First releaes' \
 -F draft=false \
 -F prerelease=false \
 -F generate_release_notes=false 


 gh api \
  --method POST \
  -H "Accept: application/vnd.github+json" \
  -H "X-GitHub-Api-Version: 2022-11-28" \
  /repos/godarkthrt/cmdHelperTools/releases/113368319/assets?name=fc-mac.zip \
  -f '@fc-map.zip'