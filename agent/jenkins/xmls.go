package jenkins

// TODO add variable replacement for dynamic paramater sent by the GUI.

// XMLCreate a job config string for jenkins.
const XMLCreate = `<?xml version='1.1' encoding='UTF-8'?>
<org.jenkinsci.plugins.workflow.multibranch.WorkflowMultiBranchProject plugin="workflow-multibranch@2.20">
  <actions/>
  <description>%v</description>
  <properties>
    <io.jenkins.blueocean.rest.impl.pipeline.credential.BlueOceanCredentialsProvider_-FolderPropertyImpl plugin="blueocean-pipeline-scm-api@1.10.1">
      <domain plugin="credentials@2.1.18">
        <name>blueocean-folder-credential-domain</name>
        <description>Blue Ocean Folder Credentials domain</description>
        <specifications/>
      </domain>
      <user>admin</user>
      <id>github-enterprise:e772171560f3215c51ebceb13f1cac470a701ac2b3463bf9737284ecc288e1ce</id>
    </io.jenkins.blueocean.rest.impl.pipeline.credential.BlueOceanCredentialsProvider_-FolderPropertyImpl>
    <org.jenkinsci.plugins.configfiles.folder.FolderConfigFileProperty plugin="config-file-provider@3.4.1">
      <configs class="sorted-set">
        <comparator class="org.jenkinsci.plugins.configfiles.folder.FolderConfigFileProperty$1"/>
      </configs>
    </org.jenkinsci.plugins.configfiles.folder.FolderConfigFileProperty>
  </properties>
  <folderViews class="jenkins.branch.MultiBranchProjectViewHolder" plugin="branch-api@2.1.2">
    <owner class="org.jenkinsci.plugins.workflow.multibranch.WorkflowMultiBranchProject" reference="../.."/>
  </folderViews>
  <healthMetrics>
    <com.cloudbees.hudson.plugins.folder.health.WorstChildHealthMetric plugin="cloudbees-folder@6.7">
      <nonRecursive>false</nonRecursive>
    </com.cloudbees.hudson.plugins.folder.health.WorstChildHealthMetric>
  </healthMetrics>
  <icon class="jenkins.branch.MetadataActionFolderIcon" plugin="branch-api@2.1.2">
    <owner class="org.jenkinsci.plugins.workflow.multibranch.WorkflowMultiBranchProject" reference="../.."/>
  </icon>
  <orphanedItemStrategy class="com.cloudbees.hudson.plugins.folder.computed.DefaultOrphanedItemStrategy" plugin="cloudbees-folder@6.7">
    <pruneDeadBranches>true</pruneDeadBranches>
    <daysToKeep>-1</daysToKeep>
    <numToKeep>-1</numToKeep>
  </orphanedItemStrategy>
  <triggers/>
  <disabled>false</disabled>
  <sources class="jenkins.branch.MultiBranchProject$BranchSourceList" plugin="branch-api@2.1.2">
    <data>
      <jenkins.branch.BranchSource>
        <source class="org.jenkinsci.plugins.github_branch_source.GitHubSCMSource" plugin="github-branch-source@2.4.2">
          <id>blueocean</id>
          <apiUri>https://github.intel.com/api/v3</apiUri>
          <credentialsId>github-enterprise:e772171560f3215c51ebceb13f1cac470a701ac2b3463bf9737284ecc288e1ce</credentialsId>
          <repoOwner>kubernetes</repoOwner>
          <repository>auto-k8s</repository>
          <traits>
            <org.jenkinsci.plugins.github__branch__source.BranchDiscoveryTrait>
              <strategyId>3</strategyId>
            </org.jenkinsci.plugins.github__branch__source.BranchDiscoveryTrait>
            <org.jenkinsci.plugins.github__branch__source.ForkPullRequestDiscoveryTrait>
              <strategyId>1</strategyId>
              <trust class="org.jenkinsci.plugins.github_branch_source.ForkPullRequestDiscoveryTrait$TrustPermission"/>
            </org.jenkinsci.plugins.github__branch__source.ForkPullRequestDiscoveryTrait>
            <org.jenkinsci.plugins.github__branch__source.OriginPullRequestDiscoveryTrait>
              <strategyId>1</strategyId>
            </org.jenkinsci.plugins.github__branch__source.OriginPullRequestDiscoveryTrait>
            <jenkins.plugins.git.traits.CleanBeforeCheckoutTrait plugin="git@3.9.1">
              <extension class="hudson.plugins.git.extensions.impl.CleanBeforeCheckout"/>
            </jenkins.plugins.git.traits.CleanBeforeCheckoutTrait>
            <jenkins.plugins.git.traits.CleanAfterCheckoutTrait plugin="git@3.9.1">
              <extension class="hudson.plugins.git.extensions.impl.CleanCheckout"/>
            </jenkins.plugins.git.traits.CleanAfterCheckoutTrait>
            <jenkins.plugins.git.traits.LocalBranchTrait plugin="git@3.9.1">
              <extension class="hudson.plugins.git.extensions.impl.LocalBranch">
                <localBranch>**</localBranch>
              </extension>
            </jenkins.plugins.git.traits.LocalBranchTrait>
          </traits>
        </source>
      </jenkins.branch.BranchSource>
    </data>
    <owner class="org.jenkinsci.plugins.workflow.multibranch.WorkflowMultiBranchProject" reference="../.."/>
  </sources>
  <factory class="org.jenkinsci.plugins.workflow.multibranch.WorkflowBranchProjectFactory">
    <owner class="org.jenkinsci.plugins.workflow.multibranch.WorkflowMultiBranchProject" reference="../.."/>
    <scriptPath>Jenkinsfile</scriptPath>
  </factory>
</org.jenkinsci.plugins.workflow.multibranch.WorkflowMultiBranchProject>`
