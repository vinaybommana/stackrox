#!/usr/bin/env bash
set -euo pipefail
IFS=$'\n\t'

main() {
    cd "$(dirname "${BASH_SOURCE[0]}")"

    echo "This script can optionally push imported images into a private registry."
    echo "Most users add a path segment like \"/stackrox\"."
    echo "For example, you might input: my-registry.example.com:5000/stackrox"
    echo "To skip pushing, simply do not enter a prefix."
    echo -n "Enter your private registry prefix: "
    read registry_prefix
    echo

    echo "Loading main image..."
    main_tag="$(docker load -i main.img | tag)"
    main_image_local="stackrox.io/main:${main_tag}"
    main_image_remote="${registry_prefix}/main:${main_tag}"

    echo "Loading monitoring image..."
    monitoring_tag="$(docker load -i monitoring.img | tag)"
    monitoring_image_local="stackrox.io/monitoring:${monitoring_tag}"
    monitoring_image_remote="${registry_prefix}/monitoring:${monitoring_tag}"

    echo "Loading scanner image..."
    scanner_tag="$(docker load -i scanner.img | tag)"
    scanner_image_local="stackrox.io/scanner:${scanner_tag}"
    scanner_image_remote="${registry_prefix}/scanner:${scanner_tag}"

    echo "Loading scanner v2 images..."
    scanner_v2_tag="$(docker load -i scanner_v2.img | tag)"
    scanner_v2_image_local="stackrox.io/scanner-v2:${scanner_v2_tag}"
    scanner_v2_image_remote="${registry_prefix}/scanner-v2:${scanner_v2_tag}"

    scanner_v2_db_tag="$(docker load -i scanner_v2_db.img | tag)"
    scanner_v2_db_image_local="stackrox.io/scanner-v2-db:${scanner_v2_db_tag}"
    scanner_v2_db_image_remote="${registry_prefix}/scanner-v2-db:${scanner_v2_db_tag}"

    if [[ -z "$registry_prefix" ]]; then
        echo "No registry prefix given, all done!"
        return
    fi

    echo "Pushing image: ${main_image_local} as ${main_image_remote}"
    docker tag "${main_image_local}" "${main_image_remote}"
    docker push "${main_image_remote}" | cat

    echo "Pushing image: ${monitoring_image_local} as ${monitoring_image_remote}"
    docker tag "${monitoring_image_local}" "${monitoring_image_remote}"
    docker push "${monitoring_image_remote}" | cat

    echo "Pushing image: ${scanner_image_local} as ${scanner_image_remote}"
    docker tag "${scanner_image_local}" "${scanner_image_remote}"
    docker push "${scanner_image_remote}" | cat

    echo "Pushing image: ${scanner_v2_image_local} as ${scanner_v2_image_remote}"
    docker tag "${scanner_v2_image_local}" "${scanner_v2_image_remote}"
    docker push "${scanner_v2_image_remote}" | cat

    echo "Pushing image: ${scanner_v2_db_image_local} as ${scanner_v2_db_image_remote}"
    docker tag "${scanner_v2_db_image_local}" "${scanner_v2_db_image_remote}"
    docker push "${scanner_v2_db_image_remote}" | cat

    echo "All done!"
}

tag() {
    sed 's/.*:\(.*$\)/\1/'
}

main
