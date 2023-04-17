#!/usr/bin/env bats

load "../helpers.bash"

setup_file() {
  # remove binaries from the previous runs
  [[ -n "$NO_BATS_ROXCTL_REBUILD" ]] || rm -f "${tmp_roxctl}"/roxctl*
  echo "Testing roxctl version: '$(roxctl-release version)'" >&3
}

setup() {
    echo "Running global flags tests"
}

teardown() {
    echo "Finished running global flags tests"
}

@test "roxctl central whoami --help should have multiple global flags" {
    run roxctl-release central whoami --help
    assert_success
    assert_multiple_global_flags
}

@test "roxctl cluster delete --help should have multiple global flags" {
    run roxctl-release cluster delete --help
    assert_success
    assert_multiple_global_flags
}

@test "roxctl collector support-packages --help should have multiple global flags" {
    run roxctl-release collector support-packages --help
    assert_success
    assert_multiple_global_flags
}

@test "roxctl completion --help should have one global flag" {
    run roxctl-release completion --help
    assert_success
    assert_single_global_flag
}

@test "roxctl deployment check --help should have multiple global flags" {
    run roxctl-release deployment check --help
    assert_success
    assert_multiple_global_flags
}

@test "roxctl generate --help should have one global flag" {
    run roxctl-release generate --help
    assert_success
    assert_single_global_flag
}

@test "roxctl helm --help should have one global flag" {
    run roxctl-release helm --help
    assert_success
    assert_single_global_flag
}

@test "roxctl image check --help should have multiple global flags" {
    run roxctl-release image check --help
    assert_success
    assert_multiple_global_flags
}

@test "roxctl scanner upload-db --help should have multiple global flags" {
    run roxctl-release scanner upload-db --help
    assert_success
    assert_multiple_global_flags
}

@test "roxctl sensor get-bundle --help should have multiple global flags" {
    run roxctl-release sensor get-bundle --help
    assert_success
    assert_multiple_global_flags
}

@test "roxctl version --help should have one global flag" {
    run roxctl-release version --help
    assert_success
    assert_single_global_flag
}
