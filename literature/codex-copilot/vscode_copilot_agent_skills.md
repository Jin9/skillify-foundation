<!DOCTYPE html>
<html lang="en">

<head>
	<meta charset="UTF-8" />
	<meta name="awa-expId" content="vscw_aaflight1016_treatment:103440;" />
	<meta name="awa-env" content="prod" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<meta name="google-site-verification" content="hNs7DXrTySP_X-0P_AC0WulAXvUwgSXEmgfcO2r79dw" />

	<!-- Twitter and Facebook OpenGraph Metadata-->
	<meta name="twitter:card" content="summary_large_image" />
	<meta name="twitter:site" content="@code" />

	<meta name="description" content="Learn how to use Agent Skills in VS Code to teach GitHub Copilot specialized capabilities that work across VS Code, GitHub Copilot CLI, and GitHub Copilot cloud agent." />
<meta name="keywords" content="" />
<meta name="ms.prod" content="vs-code" />
<meta name="ms.TOCTitle" content="" />
<meta name="ms.ContentId" content="a7d3e5f8-2c4b-4d9a-b8e1-3f6c9a2d7e41" />
<meta name="ms.date" content="4/22/2026" />
<meta name="ms.topic" content="conceptual" />
<!-- Twitter and Facebook OpenGraph Metadata-->
<meta name="twitter:card" content="summary_large_image" />
<meta property="og:url" content="https://code.visualstudio.com/docs/copilot/customization/agent-skills" />
<meta property="og:type" content="article" />
<meta property="og:title" content="Use Agent Skills in VS Code" />
<meta property="og:description" content="Learn how to use Agent Skills in VS Code to teach GitHub Copilot specialized capabilities that work across VS Code, GitHub Copilot CLI, and GitHub Copilot cloud agent." />

<meta property="og:image" content="https://code.visualstudio.com/opengraphimg/generated/docs/copilot/customization/agent-skills.webp" />



	<link rel="shortcut icon" href="/assets/favicon.ico" sizes="128x128" />
	<link rel="apple-touch-icon" href="/assets/apple-touch-icon.png">

	<title>Use Agent Skills in VS Code</title>

	<link rel="stylesheet" href="/dist/style.css">

	<script src="https://consentdeliveryfd.azurefd.net/mscc/lib/v2/wcp-consent.js"></script>
	<script type="text/javascript" src="https://js.monitor.azure.com/scripts/c/ms.analytics-web-4.min.js"></script>
	
	<script type="text/javascript">
	// Leave as var; siteConsent is initialized and referenced elsewhere.
	var siteConsent = null;
	
	const GPC_DataSharingOptIn = false;
	WcpConsent.onInitCallback(function () {
		window.appInsights = new oneDS.ApplicationInsights();
		window.appInsights.initialize({
			instrumentationKey: "1a3eb3104447440391ad5f2a6ee06a0a-62879566-bc58-4741-9650-302bf2af703f-7103",
			propertyConfiguration: {
				userConsented: false,
				gpcDataSharingOptIn: false,
				callback: {
					userConsentDetails: siteConsent ? siteConsent.getConsent : undefined
				},
			},
			cookieCfg: {
				ignoreCookies: ["MSCC"]
			},
			webAnalyticsConfiguration:{ // Web Analytics Plugin configuration
				urlCollectQuery: true,
				urlCollectHash: true,
				autoCapture: {
					scroll: true,
					pageView: true,
					onLoad: true,
					onUnload: true,
					click: true,
					resize: true,
					jsError: true
				}
			}
		}, []);
	
		window.appInsights.getPropertyManager().getPropertiesContext().web.gpcDataSharingOptIn = GPC_DataSharingOptIn;
	});
	</script>
	<link rel="alternate" type="application/atom+xml" title="RSS Feed for code.visualstudio.com" href="/feed.xml" />
</head>

<body >
	<!-- Setting theme here to avoid FOUC -->
	<script>
		function setTheme(themeName) {
			if (themeName === 'dark') {
				document.documentElement.removeAttribute('data-theme'); // dark is default, so no data-theme attribute needed
			}

			if (themeName === 'light') {
				document.documentElement.setAttribute('data-theme', themeName);
			}
			return;
		}

		// Determine initial theme: user preference or system preference
		let theme = localStorage.getItem('theme') || (window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light');
		setTheme(theme); // Apply the initial theme

		// Listen for changes in the system theme preference
		window.matchMedia('(prefers-color-scheme: dark)').addListener(e => {
			if (!localStorage.getItem('theme')) { // Only if no user preference is saved
				setTheme(e.matches ? 'dark' : 'light');
			}
		});
	</script>

	<div id="main">
		<div class="navbar-fixed-container">
			<div class="navbar navbar-inverse navbar-fixed-top ">
				<div id='cookie-banner'></div>		<nav role="navigation" aria-label="Top Level">
					<div class="container">
						<div class="nav navbar-header">
							<a class="navbar-brand" href="/"><span>Visual Studio Code</span></a>
						</div>
						<div class="navbar-collapse collapse">
							<ul class="nav navbar-nav navbar-left">
								<li class="active" ><a id="nav-docs" href="/docs">Docs</a></li>
								<li ><a id="nav-updates" href="/updates">Updates</a>
								</li>
								<li ><a id="nav-blogs" href="/blogs">Blog</a></li>
								<li ><a id="nav-extend" href="/api">API</a></li>
								<li><a href="https://marketplace.visualstudio.com/VSCode" target="_blank" rel="noopener"
										id="nav-extensions">Extensions</a></li>
								<li ><a id="nav-mcp" href="/mcp">MCP</a></li>
								<li ><a id="nav-faqs" href="/docs/supporting/faq">FAQ</a>
								</li>
								<li ><a id="nav-learn" href="/learn">Learn</a></li>
								<li><a id="nav-events" href="https://aka.ms/vscode/live" target="_blank" rel="noopener">Events</a></li>
							</ul>
							<ul class="nav navbar-nav navbar-right" role="presentation">
								<li>
									<a class="link-button" href="/Download" id="nav-download">
										<span>Download</span>
									</a>
								</li>
							</ul>
						</div>
						<div class="navbar-actions">
							<div class="search" role="presentation">
								<div class="nav-search search-control" role="button" tabindex="0" aria-label="Open search dialog">
									<div class="input-group" role="presentation">
										<span class="input-group-btn">
											<span class="btn search-icon-container" aria-hidden="true">
												<img class="search-icon-dark" src="/assets/icons/search-dark.svg" alt="" />
												<img class="search-icon-light" src="/assets/icons/search.svg" alt="" />
											</span>
										</span>
										<span class="search-box form-control" aria-hidden="true" role="presentation"></span>
										<span class="search-shortcut-placeholder">Search</span>
										<span class="search-shortcut-overlay"></span>
									</div>
								</div>					</div>
							<button type="button" class="theme-switch" id="theme-toggle">
								<img class="theme-icon-light" src="/assets/icons/theme-light.svg" alt="Switch to the dark theme" />
								<img class="theme-icon-dark" src="/assets/icons/theme-dark.svg" alt="Switch to the light theme" />
							</button>
							<a class="link-button navbar-actions-download" href="/Download">
								<span>Download</span>
							</a>
						</div>
						<button type="button" class="navbar-toggle" data-toggle="collapse" data-target=".navbar-collapse"
							aria-label="Expand and Collapse Menu">
							<span class="icon-bar"></span>
							<span class="icon-bar"></span>
							<span class="icon-bar"></span>
						</button>
					</div>
				</nav>
			</div>
		</div>		<div data-announcement-version="2026-03-17-ghcp-events" class="updates-banner js-hidden  ">
			<div class="container">
				<p class="message">Explore Agentic Development - <a href="https://aka.ms/githubcopilotdevdays" target="_blank" rel="noopener" id="banner-link-updates">Join a GitHub Copilot Dev Day near you!</a></p>
			</div>
			<div tabindex="0" role="button" title="Dismiss this update" class="dismiss-btn" id="banner-dismiss-btn"><span class="sr-only">Dismiss this update</span><span aria-hidden="true" class="glyph-icon"></span></div>
		</div>
		<!-- This div wraps around the entire site -->
		<!-- The body itself should already have a main tag -->
		<main id="main-content">
			<script>
    function closeReportIssue() {
        var element = document.getElementById('surveypopup');
        element.parentElement.removeChild(element);
    }

    function reportIssue(tutorial, page) {
        var div = document.createElement('div');
        div.innerHTML = '<div id="surveypopup" class="overlay visible"><div class="surveypopup"><div id="surveytitle">Tell us more<a href="javascript:void(0)" onclick="closeReportIssue()">X</a></div><div id="surveydiv"><iframe frameBorder="0" scrolling="0" src="https://www.research.net/r/PWZWZ52?tutorial=' + tutorial + '&step=' + page + '"></iframe></div></div></div>';
        document.body.appendChild(div.children[0]);
    }
</script>

<div class="body-content docs docs-github-layout">
    <div class="docs-layout-wrapper">
        <!-- Left sidebar - Table of Contents -->
        <aside class="docs-left-sidebar">
            <nav id="docs-navbar" aria-label="Topics" class="docs-nav visible-md visible-lg">
              <h4>Documentation</h4>
              <ul class="nav" id="main-nav">
              <li >
                <a href="/docs" >Overview</a>
              </li>
              
            <li class="panel collapsed">
              <a class="area" role="button" href="#setup-articles" data-parent="#main-nav" data-toggle="collapse">Setup</a>
              <ul id="setup-articles" class="collapse ">
            
                    <li >
                      <a href="/docs/setup/setup-overview" >Overview</a>
                    </li>
                      
                    <li >
                      <a href="/docs/setup/linux" >Linux</a>
                    </li>
                      
                    <li >
                      <a href="/docs/setup/mac" >macOS</a>
                    </li>
                      
                    <li >
                      <a href="/docs/setup/windows" >Windows</a>
                    </li>
                      
                    <li >
                      <a href="/docs/setup/vscode-web" >VS Code for the Web</a>
                    </li>
                      
                    <li >
                      <a href="/docs/setup/raspberry-pi" >Raspberry Pi</a>
                    </li>
                      
                    <li >
                      <a href="/docs/setup/network" >Network</a>
                    </li>
                      
                    <li >
                      <a href="/docs/setup/additional-components" >Additional Components</a>
                    </li>
                      
                    <li >
                      <a href="/docs/setup/uninstall" >Uninstall</a>
                    </li>
                      
              </ul>
            </li>
                
            <li class="panel collapsed">
              <a class="area" role="button" href="#getstarted-articles" data-parent="#main-nav" data-toggle="collapse">Get Started</a>
              <ul id="getstarted-articles" class="collapse ">
            
                    <li >
                      <a href="/docs/getstarted/getting-started" >VS Code Tutorial</a>
                    </li>
                      
                    <li >
                      <a href="/docs/getstarted/copilot-quickstart" >Copilot Quickstart</a>
                    </li>
                      
                    <li >
                      <a href="/docs/getstarted/userinterface" >User Interface</a>
                    </li>
                      
                    <li >
                      <a href="/docs/getstarted/personalize-vscode" >Personalize VS Code</a>
                    </li>
                      
                    <li >
                      <a href="/docs/getstarted/extensions" >Install Extensions</a>
                    </li>
                      
                    <li >
                      <a href="/docs/getstarted/tips-and-tricks" >Tips and Tricks</a>
                    </li>
                      
                    <li >
                      <a href="/docs/getstarted/introvideos" >Intro Videos</a>
                    </li>
                      
              </ul>
            </li>
                
            <li class="panel active expanded">
              <a class="area" role="button" href="#copilot-articles" data-parent="#main-nav" data-toggle="collapse">GitHub Copilot</a>
              <ul id="copilot-articles" class="collapse in">
            
                    <li >
                      <a href="/docs/copilot/overview" >Overview</a>
                    </li>
                      
                    <li >
                      <a href="/docs/copilot/setup" >Setup</a>
                    </li>
                      
                    <li >
                      <a href="/docs/copilot/getting-started" >Quickstart</a>
                    </li>
                      
            <li class="panel collapsed">
              <a class="area" role="button" href="#copilot-concepts-articles"  data-toggle="collapse">Concepts</a>
              <ul id="copilot-concepts-articles" class="collapse ">
            
                    <li >
                      <a href="/docs/copilot/concepts/overview" >Overview</a>
                    </li>
                      
                    <li >
                      <a href="/docs/copilot/concepts/language-models" >Language Models</a>
                    </li>
                      
                    <li >
                      <a href="/docs/copilot/concepts/context" >Context</a>
                    </li>
                      
                    <li >
                      <a href="/docs/copilot/concepts/tools" >Tools</a>
                    </li>
                      
                    <li >
                      <a href="/docs/copilot/concepts/agents" >Agents</a>
                    </li>
                      
                    <li >
                      <a href="/docs/copilot/concepts/customization" >Customization</a>
                    </li>
                      
                    <li >
                      <a href="/docs/copilot/concepts/trust-and-safety" >Trust & Safety</a>
                    </li>
                      
              </ul>
            </li>
                    
            <li class="panel collapsed">
              <a class="area" role="button" href="#copilot-agents-articles"  data-toggle="collapse">Agents</a>
              <ul id="copilot-agents-articles" class="collapse ">
            
                    <li >
                      <a href="/docs/copilot/agents/overview" >Overview</a>
                    </li>
                      
                    <li >
                      <a href="/docs/copilot/agents/agents-tutorial" >Agents Tutorial</a>
                    </li>
                      
                    <li >
                      <a href="/docs/copilot/agents/planning" >Planning</a>
                    </li>
                      
                    <li >
                      <a href="/docs/copilot/agents/memory" >Memory</a>
                    </li>
                      
                    <li >
                      <a href="/docs/copilot/agents/agent-tools" >Tools</a>
                    </li>
                      
                    <li >
                      <a href="/docs/copilot/agents/subagents" >Subagents</a>
                    </li>
                      
                    <li >
                      <a href="/docs/copilot/agents/local-agents" >Local Agents</a>
                    </li>
                      
                    <li >
                      <a href="/docs/copilot/agents/copilot-cli" >Copilot CLI</a>
                    </li>
                      
                    <li >
                      <a href="/docs/copilot/agents/cloud-agents" >Cloud Agents</a>
                    </li>
                      
                    <li >
                      <a href="/docs/copilot/agents/third-party-agents" >Third-Party Agents</a>
                    </li>
                      
              </ul>
            </li>
                    
            <li class="panel collapsed">
              <a class="area" role="button" href="#copilot-chat-articles"  data-toggle="collapse">Chat</a>
              <ul id="copilot-chat-articles" class="collapse ">
            
                    <li >
                      <a href="/docs/copilot/chat/copilot-chat" >Overview</a>
                    </li>
                      
                    <li >
                      <a href="/docs/copilot/chat/chat-sessions" >Chat Sessions</a>
                    </li>
                      
                    <li >
                      <a href="/docs/copilot/chat/copilot-chat-context" >Add Context</a>
                    </li>
                      
                    <li >
                      <a href="/docs/copilot/chat/inline-chat" >Inline Chat</a>
                    </li>
                      
                    <li >
                      <a href="/docs/copilot/chat/review-code-edits" >Review Edits</a>
                    </li>
                      
                    <li >
                      <a href="/docs/copilot/chat/chat-checkpoints" >Checkpoints</a>
                    </li>
                      
                    <li >
                      <a href="/docs/copilot/chat/chat-artifacts" >Artifacts Panel</a>
                    </li>
                      
                    <li >
                      <a href="/docs/copilot/chat/chat-debug-view" >Debug Chat Interactions</a>
                    </li>
                      
                    <li >
                      <a href="/docs/copilot/chat/prompt-examples" >Prompt Examples</a>
                    </li>
                      
              </ul>
            </li>
                    
            <li class="panel expanded">
              <a class="area" role="button" href="#copilot-customization-articles"  data-toggle="collapse">Customization</a>
              <ul id="copilot-customization-articles" class="collapse in">
            
                    <li >
                      <a href="/docs/copilot/customization/overview" >Overview</a>
                    </li>
                      
                    <li >
                      <a href="/docs/copilot/customization/custom-instructions" >Instructions</a>
                    </li>
                      
                    <li >
                      <a href="/docs/copilot/customization/prompt-files" >Prompt Files</a>
                    </li>
                      
                    <li >
                      <a href="/docs/copilot/customization/custom-agents" >Custom Agents</a>
                    </li>
                      
                    <li class="active">
                      <a href="/docs/copilot/customization/agent-skills" aria-label="Current Page: Agent Skills">Agent Skills</a>
                    </li>
                      
                    <li >
                      <a href="/docs/copilot/customization/language-models" >Language Models</a>
                    </li>
                      
                    <li >
                      <a href="/docs/copilot/customization/mcp-servers" >MCP</a>
                    </li>
                      
                    <li >
                      <a href="/docs/copilot/customization/hooks" >Hooks</a>
                    </li>
                      
                    <li >
                      <a href="/docs/copilot/customization/agent-plugins" >Plugins</a>
                    </li>
                      
              </ul>
            </li>
                    
            <li class="panel collapsed">
              <a class="area" role="button" href="#copilot-guides-articles"  data-toggle="collapse">Guides & Tutorials</a>
              <ul id="copilot-guides-articles" class="collapse ">
            
                    <li >
                      <a href="/docs/copilot/guides/context-engineering-guide" >Context Engineering</a>
                    </li>
                      
                    <li >
                      <a href="/docs/copilot/guides/customize-copilot-guide" >Customize AI</a>
                    </li>
                      
                    <li >
                      <a href="/docs/copilot/guides/test-driven-development-guide" >Test-Driven Development</a>
                    </li>
                      
                    <li >
                      <a href="/docs/copilot/guides/notebooks-with-ai" >Edit Notebooks with AI</a>
                    </li>
                      
                    <li >
                      <a href="/docs/copilot/guides/test-with-copilot" >Test with AI</a>
                    </li>
                      
                    <li >
                      <a href="/docs/copilot/guides/browser-agent-testing-guide" >Test Web Apps with Browser Tools</a>
                    </li>
                      
                    <li >
                      <a href="/docs/copilot/guides/debug-with-copilot" >Debug with AI</a>
                    </li>
                      
                    <li >
                      <a href="/docs/copilot/guides/mcp-developer-guide" >MCP Dev Guide</a>
                    </li>
                      
                    <li >
                      <a href="/docs/copilot/guides/monitoring-agents" >Monitoring</a>
                    </li>
                      
              </ul>
            </li>
                    
                    <li >
                      <a href="/docs/copilot/ai-powered-suggestions" >Inline Suggestions</a>
                    </li>
                      
                    <li >
                      <a href="/docs/copilot/copilot-smart-actions" >Smart Actions</a>
                    </li>
                      
                    <li >
                      <a href="/docs/copilot/best-practices" >Best Practices</a>
                    </li>
                      
                    <li >
                      <a href="/docs/copilot/security" >Security</a>
                    </li>
                      
                    <li >
                      <a href="/docs/copilot/troubleshooting" >Troubleshooting</a>
                    </li>
                      
                    <li >
                      <a href="/docs/copilot/faq" >FAQ</a>
                    </li>
                      
            <li class="panel collapsed">
              <a class="area" role="button" href="#copilot-reference-articles"  data-toggle="collapse">Reference</a>
              <ul id="copilot-reference-articles" class="collapse ">
            
                    <li >
                      <a href="/docs/copilot/reference/copilot-vscode-features" >Cheat Sheet</a>
                    </li>
                      
                    <li >
                      <a href="/docs/copilot/reference/copilot-settings" >Settings Reference</a>
                    </li>
                      
                    <li >
                      <a href="/docs/copilot/reference/mcp-configuration" >MCP Configuration</a>
                    </li>
                      
                    <li >
                      <a href="/docs/copilot/reference/workspace-context" >Workspace Context</a>
                    </li>
                      
              </ul>
            </li>
                    
              </ul>
            </li>
                
            <li class="panel collapsed">
              <a class="area" role="button" href="#configure-articles" data-parent="#main-nav" data-toggle="collapse">Configure</a>
              <ul id="configure-articles" class="collapse ">
            
                    <li >
                      <a href="/docs/configure/locales" >Display Language</a>
                    </li>
                      
                    <li >
                      <a href="/docs/configure/custom-layout" >Layout</a>
                    </li>
                      
                    <li >
                      <a href="/docs/configure/keybindings" >Keyboard Shortcuts</a>
                    </li>
                      
                    <li >
                      <a href="/docs/configure/settings" >Settings</a>
                    </li>
                      
                    <li >
                      <a href="/docs/configure/settings-sync" >Settings Sync</a>
                    </li>
                      
            <li class="panel collapsed">
              <a class="area" role="button" href="#configure-extensions-articles"  data-toggle="collapse">Extensions</a>
              <ul id="configure-extensions-articles" class="collapse ">
            
                    <li >
                      <a href="/docs/configure/extensions/extension-marketplace" >Extension Marketplace</a>
                    </li>
                      
                    <li >
                      <a href="/docs/configure/extensions/extension-runtime-security" >Extension Runtime Security</a>
                    </li>
                      
              </ul>
            </li>
                    
                    <li >
                      <a href="/docs/configure/themes" >Themes</a>
                    </li>
                      
                    <li >
                      <a href="/docs/configure/profiles" >Profiles</a>
                    </li>
                      
            <li class="panel collapsed">
              <a class="area" role="button" href="#configure-accessibility-articles"  data-toggle="collapse">Accessibility</a>
              <ul id="configure-accessibility-articles" class="collapse ">
            
                    <li >
                      <a href="/docs/configure/accessibility/accessibility" >Overview</a>
                    </li>
                      
                    <li >
                      <a href="/docs/configure/accessibility/voice" >Voice Interactions</a>
                    </li>
                      
              </ul>
            </li>
                    
                    <li >
                      <a href="/docs/configure/command-line" >Command Line Interface</a>
                    </li>
                      
                    <li >
                      <a href="/docs/configure/telemetry" >Telemetry</a>
                    </li>
                      
              </ul>
            </li>
                
            <li class="panel collapsed">
              <a class="area" role="button" href="#editing-articles" data-parent="#main-nav" data-toggle="collapse">Edit Code</a>
              <ul id="editing-articles" class="collapse ">
            
                    <li >
                      <a href="/docs/editing/codebasics" >Basic Editing</a>
                    </li>
                      
                    <li >
                      <a href="/docs/editing/intellisense" >IntelliSense</a>
                    </li>
                      
                    <li >
                      <a href="/docs/editing/editingevolved" >Code Navigation</a>
                    </li>
                      
                    <li >
                      <a href="/docs/editing/refactoring" >Refactoring</a>
                    </li>
                      
                    <li >
                      <a href="/docs/editing/userdefinedsnippets" >Snippets</a>
                    </li>
                      
            <li class="panel collapsed">
              <a class="area" role="button" href="#editing-workspaces-articles"  data-toggle="collapse">Workspaces</a>
              <ul id="editing-workspaces-articles" class="collapse ">
            
                    <li >
                      <a href="/docs/editing/workspaces/workspaces" >Overview</a>
                    </li>
                      
                    <li >
                      <a href="/docs/editing/workspaces/multi-root-workspaces" >Multi-Root Workspaces</a>
                    </li>
                      
                    <li >
                      <a href="/docs/editing/workspaces/workspace-trust" >Workspace Trust</a>
                    </li>
                      
              </ul>
            </li>
                    
              </ul>
            </li>
                
            <li class="panel collapsed">
              <a class="area" role="button" href="#debugtest-articles" data-parent="#main-nav" data-toggle="collapse">Build, Debug, Test</a>
              <ul id="debugtest-articles" class="collapse ">
            
                    <li >
                      <a href="/docs/debugtest/tasks" >Tasks</a>
                    </li>
                      
                    <li >
                      <a href="/docs/debugtest/debugging" >Debugging</a>
                    </li>
                      
                    <li >
                      <a href="/docs/debugtest/debugging-configuration" >Debug Configuration</a>
                    </li>
                      
                    <li >
                      <a href="/docs/debugtest/testing" >Testing</a>
                    </li>
                      
                    <li >
                      <a href="/docs/debugtest/port-forwarding" >Port Forwarding</a>
                    </li>
                      
                    <li >
                      <a href="/docs/debugtest/integrated-browser" >Integrated Browser</a>
                    </li>
                      
              </ul>
            </li>
                
            <li class="panel collapsed">
              <a class="area" role="button" href="#sourcecontrol-articles" data-parent="#main-nav" data-toggle="collapse">Source Control</a>
              <ul id="sourcecontrol-articles" class="collapse ">
            
                    <li >
                      <a href="/docs/sourcecontrol/overview" >Overview</a>
                    </li>
                      
                    <li >
                      <a href="/docs/sourcecontrol/quickstart" >Quickstart</a>
                    </li>
                      
                    <li >
                      <a href="/docs/sourcecontrol/staging-commits" >Staging & Committing</a>
                    </li>
                      
                    <li >
                      <a href="/docs/sourcecontrol/branches-worktrees" >Branches & Worktrees</a>
                    </li>
                      
                    <li >
                      <a href="/docs/sourcecontrol/repos-remotes" >Repositories & Remotes</a>
                    </li>
                      
                    <li >
                      <a href="/docs/sourcecontrol/merge-conflicts" >Merge Conflicts</a>
                    </li>
                      
                    <li >
                      <a href="/docs/sourcecontrol/github" >Collaborate on GitHub</a>
                    </li>
                      
                    <li >
                      <a href="/docs/sourcecontrol/troubleshooting" >Troubleshooting</a>
                    </li>
                      
                    <li >
                      <a href="/docs/sourcecontrol/faq" >FAQ</a>
                    </li>
                      
              </ul>
            </li>
                
            <li class="panel collapsed">
              <a class="area" role="button" href="#terminal-articles" data-parent="#main-nav" data-toggle="collapse">Terminal</a>
              <ul id="terminal-articles" class="collapse ">
            
                    <li >
                      <a href="/docs/terminal/getting-started" >Getting Started Tutorial</a>
                    </li>
                      
                    <li >
                      <a href="/docs/terminal/basics" >Terminal Basics</a>
                    </li>
                      
                    <li >
                      <a href="/docs/terminal/profiles" >Terminal Profiles</a>
                    </li>
                      
                    <li >
                      <a href="/docs/terminal/shell-integration" >Shell Integration</a>
                    </li>
                      
                    <li >
                      <a href="/docs/terminal/appearance" >Appearance</a>
                    </li>
                      
                    <li >
                      <a href="/docs/terminal/advanced" >Advanced</a>
                    </li>
                      
              </ul>
            </li>
                
            <li class="panel collapsed">
              <a class="area" role="button" href="#enterprise-articles" data-parent="#main-nav" data-toggle="collapse">Enterprise</a>
              <ul id="enterprise-articles" class="collapse ">
            
                    <li >
                      <a href="/docs/enterprise/overview" >Overview</a>
                    </li>
                      
                    <li >
                      <a href="/docs/enterprise/policies" >Enterprise Policies</a>
                    </li>
                      
                    <li >
                      <a href="/docs/enterprise/ai-settings" >AI Settings</a>
                    </li>
                      
                    <li >
                      <a href="/docs/enterprise/extensions" >Extensions</a>
                    </li>
                      
                    <li >
                      <a href="/docs/enterprise/telemetry" >Telemetry</a>
                    </li>
                      
                    <li >
                      <a href="/docs/enterprise/updates" >Updates</a>
                    </li>
                      
              </ul>
            </li>
                
            <li class="panel collapsed">
              <a class="area" role="button" href="#languages-articles" data-parent="#main-nav" data-toggle="collapse">Languages</a>
              <ul id="languages-articles" class="collapse ">
            
                    <li >
                      <a href="/docs/languages/overview" >Overview</a>
                    </li>
                      
                    <li >
                      <a href="/docs/languages/javascript" >JavaScript</a>
                    </li>
                      
                    <li >
                      <a href="/docs/languages/json" >JSON</a>
                    </li>
                      
                    <li >
                      <a href="/docs/languages/html" >HTML</a>
                    </li>
                      
                    <li >
                      <a href="/docs/languages/emmet" >Emmet</a>
                    </li>
                      
                    <li >
                      <a href="/docs/languages/css" >CSS, SCSS and Less</a>
                    </li>
                      
                    <li >
                      <a href="/docs/languages/typescript" >TypeScript</a>
                    </li>
                      
                    <li >
                      <a href="/docs/languages/markdown" >Markdown</a>
                    </li>
                      
                    <li >
                      <a href="/docs/languages/powershell" >PowerShell</a>
                    </li>
                      
                    <li >
                      <a href="/docs/languages/cpp" >C++</a>
                    </li>
                      
                    <li >
                      <a href="/docs/languages/java" >Java</a>
                    </li>
                      
                    <li >
                      <a href="/docs/languages/php" >PHP</a>
                    </li>
                      
                    <li >
                      <a href="/docs/languages/python" >Python</a>
                    </li>
                      
                    <li >
                      <a href="/docs/languages/julia" >Julia</a>
                    </li>
                      
                    <li >
                      <a href="/docs/languages/r" >R</a>
                    </li>
                      
                    <li >
                      <a href="/docs/languages/ruby" >Ruby</a>
                    </li>
                      
                    <li >
                      <a href="/docs/languages/rust" >Rust</a>
                    </li>
                      
                    <li >
                      <a href="/docs/languages/go" >Go</a>
                    </li>
                      
                    <li >
                      <a href="/docs/languages/tsql" >T-SQL</a>
                    </li>
                      
                    <li >
                      <a href="/docs/languages/csharp" >C#</a>
                    </li>
                      
                    <li >
                      <a href="/docs/languages/dotnet" >.NET</a>
                    </li>
                      
                    <li >
                      <a href="/docs/languages/swift" >Swift</a>
                    </li>
                      
              </ul>
            </li>
                
            <li class="panel collapsed">
              <a class="area" role="button" href="#nodejs-articles" data-parent="#main-nav" data-toggle="collapse">Node.js / JavaScript</a>
              <ul id="nodejs-articles" class="collapse ">
            
                    <li >
                      <a href="/docs/nodejs/working-with-javascript" >Working with JavaScript</a>
                    </li>
                      
                    <li >
                      <a href="/docs/nodejs/nodejs-tutorial" >Node.js Tutorial</a>
                    </li>
                      
                    <li >
                      <a href="/docs/nodejs/nodejs-debugging" >Node.js Debugging</a>
                    </li>
                      
                    <li >
                      <a href="/docs/nodejs/nodejs-deployment" >Deploy Node.js Apps</a>
                    </li>
                      
                    <li >
                      <a href="/docs/nodejs/browser-debugging" >Browser Debugging</a>
                    </li>
                      
                    <li >
                      <a href="/docs/nodejs/angular-tutorial" >Angular Tutorial</a>
                    </li>
                      
                    <li >
                      <a href="/docs/nodejs/reactjs-tutorial" >React Tutorial</a>
                    </li>
                      
                    <li >
                      <a href="/docs/nodejs/vuejs-tutorial" >Vue Tutorial</a>
                    </li>
                      
                    <li >
                      <a href="/docs/nodejs/debugging-recipes" >Debugging Recipes</a>
                    </li>
                      
                    <li >
                      <a href="/docs/nodejs/profiling" >Performance Profiling</a>
                    </li>
                      
                    <li >
                      <a href="/docs/nodejs/extensions" >Extensions</a>
                    </li>
                      
              </ul>
            </li>
                
            <li class="panel collapsed">
              <a class="area" role="button" href="#typescript-articles" data-parent="#main-nav" data-toggle="collapse">TypeScript</a>
              <ul id="typescript-articles" class="collapse ">
            
                    <li >
                      <a href="/docs/typescript/typescript-tutorial" >Tutorial</a>
                    </li>
                      
                    <li >
                      <a href="/docs/typescript/typescript-transpiling" >Transpiling</a>
                    </li>
                      
                    <li >
                      <a href="/docs/typescript/typescript-editing" >Editing</a>
                    </li>
                      
                    <li >
                      <a href="/docs/typescript/typescript-refactoring" >Refactoring</a>
                    </li>
                      
                    <li >
                      <a href="/docs/typescript/typescript-debugging" >Debugging</a>
                    </li>
                      
              </ul>
            </li>
                
            <li class="panel collapsed">
              <a class="area" role="button" href="#python-articles" data-parent="#main-nav" data-toggle="collapse">Python</a>
              <ul id="python-articles" class="collapse ">
            
                    <li >
                      <a href="/docs/python/python-quick-start" >Quick Start</a>
                    </li>
                      
                    <li >
                      <a href="/docs/python/python-tutorial" >Tutorial</a>
                    </li>
                      
                    <li >
                      <a href="/docs/python/run" >Run Python Code</a>
                    </li>
                      
                    <li >
                      <a href="/docs/python/editing" >Editing</a>
                    </li>
                      
                    <li >
                      <a href="/docs/python/linting" >Linting</a>
                    </li>
                      
                    <li >
                      <a href="/docs/python/formatting" >Formatting</a>
                    </li>
                      
                    <li >
                      <a href="/docs/python/debugging" >Debugging</a>
                    </li>
                      
                    <li >
                      <a href="/docs/python/environments" >Environments</a>
                    </li>
                      
                    <li >
                      <a href="/docs/python/testing" >Testing</a>
                    </li>
                      
                    <li >
                      <a href="/docs/python/jupyter-support-py" >Python Interactive</a>
                    </li>
                      
                    <li >
                      <a href="/docs/python/tutorial-django" >Django Tutorial</a>
                    </li>
                      
                    <li >
                      <a href="/docs/python/tutorial-fastapi" >FastAPI Tutorial</a>
                    </li>
                      
                    <li >
                      <a href="/docs/python/tutorial-flask" >Flask Tutorial</a>
                    </li>
                      
                    <li >
                      <a href="/docs/python/tutorial-create-containers" >Create Containers</a>
                    </li>
                      
                    <li >
                      <a href="/docs/python/python-on-azure" >Deploy Python Apps</a>
                    </li>
                      
                    <li >
                      <a href="/docs/python/python-web" >Python in the Web</a>
                    </li>
                      
                    <li >
                      <a href="/docs/python/settings-reference" >Settings Reference</a>
                    </li>
                      
              </ul>
            </li>
                
            <li class="panel collapsed">
              <a class="area" role="button" href="#java-articles" data-parent="#main-nav" data-toggle="collapse">Java</a>
              <ul id="java-articles" class="collapse ">
            
                    <li >
                      <a href="/docs/java/java-tutorial" >Getting Started</a>
                    </li>
                      
                    <li >
                      <a href="/docs/java/java-editing" >Navigate and Edit</a>
                    </li>
                      
                    <li >
                      <a href="/docs/java/java-refactoring" >Refactoring</a>
                    </li>
                      
                    <li >
                      <a href="/docs/java/java-linting" >Formatting and Linting</a>
                    </li>
                      
                    <li >
                      <a href="/docs/java/java-project" >Project Management</a>
                    </li>
                      
                    <li >
                      <a href="/docs/java/java-build" >Build Tools</a>
                    </li>
                      
                    <li >
                      <a href="/docs/java/java-debugging" >Run and Debug</a>
                    </li>
                      
                    <li >
                      <a href="/docs/java/java-testing" >Testing</a>
                    </li>
                      
                    <li >
                      <a href="/docs/java/java-spring-boot" >Spring Boot</a>
                    </li>
                      
                    <li >
                      <a href="/docs/java/java-app-mod" >Modernizing Java Apps</a>
                    </li>
                      
                    <li >
                      <a href="/docs/java/java-tomcat-jetty" >Application Servers</a>
                    </li>
                      
                    <li >
                      <a href="/docs/java/java-on-azure" >Deploy Java Apps</a>
                    </li>
                      
                    <li >
                      <a href="/docs/java/java-gui" >GUI Applications</a>
                    </li>
                      
                    <li >
                      <a href="/docs/java/extensions" >Extensions</a>
                    </li>
                      
                    <li >
                      <a href="/docs/java/java-faq" >FAQ</a>
                    </li>
                      
              </ul>
            </li>
                
            <li class="panel collapsed">
              <a class="area" role="button" href="#cpp-articles" data-parent="#main-nav" data-toggle="collapse">C++</a>
              <ul id="cpp-articles" class="collapse ">
            
                    <li >
                      <a href="/docs/cpp/introvideos-cpp" >Intro Videos</a>
                    </li>
                      
                    <li >
                      <a href="/docs/cpp/config-linux" >GCC on Linux</a>
                    </li>
                      
                    <li >
                      <a href="/docs/cpp/config-mingw" >GCC on Windows</a>
                    </li>
                      
                    <li >
                      <a href="/docs/cpp/config-wsl" >GCC on Windows Subsystem for Linux</a>
                    </li>
                      
                    <li >
                      <a href="/docs/cpp/config-clang-mac" >Clang on macOS</a>
                    </li>
                      
                    <li >
                      <a href="/docs/cpp/config-msvc" >Microsoft C++ on Windows</a>
                    </li>
                      
                    <li >
                      <a href="/docs/cpp/build-with-cmake" >Build with CMake</a>
                    </li>
                      
                    <li >
                      <a href="/docs/cpp/cmake-linux" >CMake Tools on Linux</a>
                    </li>
                      
                    <li >
                      <a href="/docs/cpp/cmake-quickstart" >CMake Quick Start</a>
                    </li>
                      
                    <li >
                      <a href="/docs/cpp/cpp-devtools" >C++ Dev Tools for Copilot</a>
                    </li>
                      
                    <li >
                      <a href="/docs/cpp/cpp-ide" >Editing and Navigating</a>
                    </li>
                      
                    <li >
                      <a href="/docs/cpp/cpp-debug" >Debugging</a>
                    </li>
                      
                    <li >
                      <a href="/docs/cpp/launch-json-reference" >Configure Debugging</a>
                    </li>
                      
                    <li >
                      <a href="/docs/cpp/cpp-refactoring" >Refactoring</a>
                    </li>
                      
                    <li >
                      <a href="/docs/cpp/customize-cpp-settings" >Settings Reference</a>
                    </li>
                      
                    <li >
                      <a href="/docs/cpp/configure-intellisense" >Configure IntelliSense</a>
                    </li>
                      
                    <li >
                      <a href="/docs/cpp/configure-intellisense-crosscompilation" >Configure IntelliSense for Cross-Compiling</a>
                    </li>
                      
                    <li >
                      <a href="/docs/cpp/faq-cpp" >FAQ</a>
                    </li>
                      
              </ul>
            </li>
                
            <li class="panel collapsed">
              <a class="area" role="button" href="#csharp-articles" data-parent="#main-nav" data-toggle="collapse">C#</a>
              <ul id="csharp-articles" class="collapse ">
            
                    <li >
                      <a href="/docs/csharp/introvideos-csharp" >Intro Videos</a>
                    </li>
                      
                    <li >
                      <a href="/docs/csharp/get-started" >Get Started</a>
                    </li>
                      
                    <li >
                      <a href="/docs/csharp/navigate-edit" >Navigate and Edit</a>
                    </li>
                      
                    <li >
                      <a href="/docs/csharp/intellicode" >IntelliCode</a>
                    </li>
                      
                    <li >
                      <a href="/docs/csharp/refactoring" >Refactoring</a>
                    </li>
                      
                    <li >
                      <a href="/docs/csharp/formatting-linting" >Formatting and Linting</a>
                    </li>
                      
                    <li >
                      <a href="/docs/csharp/project-management" >Project Management</a>
                    </li>
                      
                    <li >
                      <a href="/docs/csharp/build-tools" >Build Tools</a>
                    </li>
                      
                    <li >
                      <a href="/docs/csharp/package-management" >Package Management</a>
                    </li>
                      
                    <li >
                      <a href="/docs/csharp/debugging" >Run and Debug</a>
                    </li>
                      
                    <li >
                      <a href="/docs/csharp/testing" >Testing</a>
                    </li>
                      
                    <li >
                      <a href="/docs/csharp/cs-dev-kit-faq" >FAQ</a>
                    </li>
                      
              </ul>
            </li>
                
            <li class="panel collapsed">
              <a class="area" role="button" href="#containers-articles" data-parent="#main-nav" data-toggle="collapse">Container Tools</a>
              <ul id="containers-articles" class="collapse ">
            
                    <li >
                      <a href="/docs/containers/overview" >Overview</a>
                    </li>
                      
                    <li >
                      <a href="/docs/containers/quickstart-node" >Node.js</a>
                    </li>
                      
                    <li >
                      <a href="/docs/containers/quickstart-python" >Python</a>
                    </li>
                      
                    <li >
                      <a href="/docs/containers/quickstart-aspnet-core" >ASP.NET Core</a>
                    </li>
                      
                    <li >
                      <a href="/docs/containers/debug-common" >Debug</a>
                    </li>
                      
                    <li >
                      <a href="/docs/containers/docker-compose" >Docker Compose</a>
                    </li>
                      
                    <li >
                      <a href="/docs/containers/quickstart-container-registries" >Registries</a>
                    </li>
                      
                    <li >
                      <a href="/docs/containers/app-service" >Deploy to Azure</a>
                    </li>
                      
                    <li >
                      <a href="/docs/containers/choosing-dev-environment" >Choose a Dev Environment</a>
                    </li>
                      
                    <li >
                      <a href="/docs/containers/reference" >Customize</a>
                    </li>
                      
                    <li >
                      <a href="/docs/containers/bridge-to-kubernetes" >Develop with Kubernetes</a>
                    </li>
                      
                    <li >
                      <a href="/docs/containers/troubleshooting" >Tips and Tricks</a>
                    </li>
                      
              </ul>
            </li>
                
            <li class="panel collapsed">
              <a class="area" role="button" href="#datascience-articles" data-parent="#main-nav" data-toggle="collapse">Data Science</a>
              <ul id="datascience-articles" class="collapse ">
            
                    <li >
                      <a href="/docs/datascience/overview" >Overview</a>
                    </li>
                      
                    <li >
                      <a href="/docs/datascience/jupyter-notebooks" >Jupyter Notebooks</a>
                    </li>
                      
                    <li >
                      <a href="/docs/datascience/data-science-tutorial" >Data Science Tutorial</a>
                    </li>
                      
                    <li >
                      <a href="/docs/datascience/python-interactive" >Python Interactive</a>
                    </li>
                      
                    <li >
                      <a href="/docs/datascience/data-wrangler-quick-start" >Data Wrangler Quick Start</a>
                    </li>
                      
                    <li >
                      <a href="/docs/datascience/data-wrangler" >Data Wrangler</a>
                    </li>
                      
                    <li >
                      <a href="/docs/datascience/pytorch-support" >PyTorch Support</a>
                    </li>
                      
                    <li >
                      <a href="/docs/datascience/azure-machine-learning" >Azure Machine Learning</a>
                    </li>
                      
                    <li >
                      <a href="/docs/datascience/jupyter-kernel-management" >Manage Jupyter Kernels</a>
                    </li>
                      
                    <li >
                      <a href="/docs/datascience/notebooks-web" >Jupyter Notebooks on the Web</a>
                    </li>
                      
                    <li >
                      <a href="/docs/datascience/microsoft-fabric-quickstart" >Data Science in Microsoft Fabric</a>
                    </li>
                      
              </ul>
            </li>
                
            <li class="panel collapsed">
              <a class="area" role="button" href="#intelligentapps-articles" data-parent="#main-nav" data-toggle="collapse">Intelligent Apps</a>
              <ul id="intelligentapps-articles" class="collapse ">
            
                    <li >
                      <a href="/docs/intelligentapps/overview" >Foundry Toolkit Overview</a>
                    </li>
                      
                    <li >
                      <a href="/docs/intelligentapps/copilot-tools" >Foundry Toolkit Copilot Tools</a>
                    </li>
                      
                    <li >
                      <a href="/docs/intelligentapps/create-agents" >Create Agents</a>
                    </li>
                      
                    <li >
                      <a href="/docs/intelligentapps/models" >Models</a>
                    </li>
                      
                    <li >
                      <a href="/docs/intelligentapps/playground" >Playground</a>
                    </li>
                      
                    <li >
                      <a href="/docs/intelligentapps/agentbuilder" >Agent Builder</a>
                    </li>
                      
                    <li >
                      <a href="/docs/intelligentapps/agent-inspector" >Agent Inspector</a>
                    </li>
                      
                    <li >
                      <a href="/docs/intelligentapps/evaluation" >Evaluation</a>
                    </li>
                      
                    <li >
                      <a href="/docs/intelligentapps/tool-catalog" >Tool Catalog</a>
                    </li>
                      
                    <li >
                      <a href="/docs/intelligentapps/finetune" >Fine-Tuning (Automated Setup)</a>
                    </li>
                      
                    <li >
                      <a href="/docs/intelligentapps/finetune-legacy" >Fine-Tuning (Project Template)</a>
                    </li>
                      
                    <li >
                      <a href="/docs/intelligentapps/modelconversion" >Model Conversion</a>
                    </li>
                      
                    <li >
                      <a href="/docs/intelligentapps/tracing" >Tracing</a>
                    </li>
                      
                    <li >
                      <a href="/docs/intelligentapps/profiling" >Profiling (Windows ML)</a>
                    </li>
                      
                    <li >
                      <a href="/docs/intelligentapps/faq" >FAQ</a>
                    </li>
                      
            <li class="panel collapsed">
              <a class="area" role="button" href="#intelligentapps-reference-articles"  data-toggle="collapse">Reference</a>
              <ul id="intelligentapps-reference-articles" class="collapse ">
            
                    <li >
                      <a href="/docs/intelligentapps/reference/FileStructure" >File Structure</a>
                    </li>
                      
                    <li >
                      <a href="/docs/intelligentapps/reference/ManualModelConversion" >Manual Model Conversion</a>
                    </li>
                      
                    <li >
                      <a href="/docs/intelligentapps/reference/ManualConversionOnGPU" >Manual Model Conversion on GPU</a>
                    </li>
                      
                    <li >
                      <a href="/docs/intelligentapps/reference/SetupWithoutAITK" >Setup Environment Without Foundry Toolkit</a>
                    </li>
                      
                    <li >
                      <a href="/docs/intelligentapps/reference/TemplateProject" >Template Project</a>
                    </li>
                      
                    <li >
                      <a href="/docs/intelligentapps/reference/migrate-from-visualizer" >Migrating from Visualizer to Agent Inspector</a>
                    </li>
                      
              </ul>
            </li>
                    
              </ul>
            </li>
                
            <li class="panel collapsed">
              <a class="area" role="button" href="#azure-articles" data-parent="#main-nav" data-toggle="collapse">Azure</a>
              <ul id="azure-articles" class="collapse ">
            
                    <li >
                      <a href="/docs/azure/overview" >Overview</a>
                    </li>
                      
                    <li >
                      <a href="/docs/azure/gettingstarted" >Getting Started</a>
                    </li>
                      
                    <li >
                      <a href="/docs/azure/resourcesextension" >Resources View</a>
                    </li>
                      
                    <li >
                      <a href="/docs/azure/deployment" >Deployment</a>
                    </li>
                      
                    <li >
                      <a href="/docs/azure/vscodeforweb" >VS Code for the Web - Azure</a>
                    </li>
                      
                    <li >
                      <a href="/docs/azure/containers" >Containers</a>
                    </li>
                      
                    <li >
                      <a href="/docs/azure/aksextensions" >Azure Kubernetes Service</a>
                    </li>
                      
                    <li >
                      <a href="/docs/azure/kubernetes" >Kubernetes</a>
                    </li>
                      
                    <li >
                      <a href="/docs/azure/mongodb" >MongoDB</a>
                    </li>
                      
                    <li >
                      <a href="/docs/azure/remote-debugging" >Remote Debugging for Node.js</a>
                    </li>
                      
              </ul>
            </li>
                
            <li class="panel collapsed">
              <a class="area" role="button" href="#remote-articles" data-parent="#main-nav" data-toggle="collapse">Remote</a>
              <ul id="remote-articles" class="collapse ">
            
                    <li >
                      <a href="/docs/remote/remote-overview" >Overview</a>
                    </li>
                      
                    <li >
                      <a href="/docs/remote/ssh" >SSH</a>
                    </li>
                      
                    <li >
                      <a href="/docs/remote/dev-containers" >Dev Containers</a>
                    </li>
                      
                    <li >
                      <a href="/docs/remote/wsl" >Windows Subsystem for Linux</a>
                    </li>
                      
                    <li >
                      <a href="/docs/remote/codespaces" >GitHub Codespaces</a>
                    </li>
                      
                    <li >
                      <a href="/docs/remote/vscode-server" >VS Code Server</a>
                    </li>
                      
                    <li >
                      <a href="/docs/remote/tunnels" >Tunnels</a>
                    </li>
                      
                    <li >
                      <a href="/docs/remote/ssh-tutorial" >SSH Tutorial</a>
                    </li>
                      
                    <li >
                      <a href="/docs/remote/wsl-tutorial" >WSL Tutorial</a>
                    </li>
                      
                    <li >
                      <a href="/docs/remote/troubleshooting" >Tips and Tricks</a>
                    </li>
                      
                    <li >
                      <a href="/docs/remote/faq" >FAQ</a>
                    </li>
                      
              </ul>
            </li>
                
            <li class="panel collapsed">
              <a class="area" role="button" href="#devcontainers-articles" data-parent="#main-nav" data-toggle="collapse">Dev Containers</a>
              <ul id="devcontainers-articles" class="collapse ">
            
                    <li >
                      <a href="/docs/devcontainers/containers" >Overview</a>
                    </li>
                      
                    <li >
                      <a href="/docs/devcontainers/tutorial" >Tutorial</a>
                    </li>
                      
                    <li >
                      <a href="/docs/devcontainers/attach-container" >Attach to Container</a>
                    </li>
                      
                    <li >
                      <a href="/docs/devcontainers/create-dev-container" >Create Dev Container</a>
                    </li>
                      
                    <li >
                      <a href="/docs/devcontainers/containers-advanced" >Advanced Containers</a>
                    </li>
                      
                    <li >
                      <a href="/docs/devcontainers/devcontainerjson-reference" >devcontainer.json</a>
                    </li>
                      
                    <li >
                      <a href="/docs/devcontainers/devcontainer-cli" >Dev Container CLI</a>
                    </li>
                      
                    <li >
                      <a href="/docs/devcontainers/tips-and-tricks" >Tips and Tricks</a>
                    </li>
                      
                    <li >
                      <a href="/docs/devcontainers/faq" >FAQ</a>
                    </li>
                      
              </ul>
            </li>
                
            <li class="panel collapsed">
              <a class="area" role="button" href="#reference-articles" data-parent="#main-nav" data-toggle="collapse">Reference</a>
              <ul id="reference-articles" class="collapse ">
            
                    <li >
                      <a href="/docs/reference/default-keybindings" >Default Keyboard Shortcuts</a>
                    </li>
                      
                    <li >
                      <a href="/docs/reference/default-settings" >Default Settings</a>
                    </li>
                      
                    <li >
                      <a href="/docs/reference/variables-reference" >Substitution Variables</a>
                    </li>
                      
                    <li >
                      <a href="/docs/reference/tasks-appendix" >Tasks Schema</a>
                    </li>
                      
              </ul>
            </li>
                
              </ul>
            </nav>
            
            <nav id="small-nav" aria-label="Topics" class="docs-nav hidden-md hidden-lg">
              <label class="faux-h4" for="small-nav-dropdown">Topics</label>
              <select id="small-nav-dropdown" aria-label="topics">
                <option value="/docs" >Overview</option>
            
              <optgroup label="Setup">
            
                    <option value="/docs/setup/setup-overview" >Overview</option>
                    
                    <option value="/docs/setup/linux" >Linux</option>
                    
                    <option value="/docs/setup/mac" >macOS</option>
                    
                    <option value="/docs/setup/windows" >Windows</option>
                    
                    <option value="/docs/setup/vscode-web" >VS Code for the Web</option>
                    
                    <option value="/docs/setup/raspberry-pi" >Raspberry Pi</option>
                    
                    <option value="/docs/setup/network" >Network</option>
                    
                    <option value="/docs/setup/additional-components" >Additional Components</option>
                    
                    <option value="/docs/setup/uninstall" >Uninstall</option>
                    
              </optgroup>
              
              <optgroup label="Get Started">
            
                    <option value="/docs/getstarted/getting-started" >VS Code Tutorial</option>
                    
                    <option value="/docs/getstarted/copilot-quickstart" >Copilot Quickstart</option>
                    
                    <option value="/docs/getstarted/userinterface" >User Interface</option>
                    
                    <option value="/docs/getstarted/personalize-vscode" >Personalize VS Code</option>
                    
                    <option value="/docs/getstarted/extensions" >Install Extensions</option>
                    
                    <option value="/docs/getstarted/tips-and-tricks" >Tips and Tricks</option>
                    
                    <option value="/docs/getstarted/introvideos" >Intro Videos</option>
                    
              </optgroup>
              
              <optgroup label="GitHub Copilot">
            
                    <option value="/docs/copilot/overview" >Overview</option>
                    
                    <option value="/docs/copilot/setup" >Setup</option>
                    
                    <option value="/docs/copilot/getting-started" >Quickstart</option>
                    
                    <optgroup label="GitHub Copilot - Concepts">
            
                      <option value="/docs/copilot/concepts/overview" >Overview</option>
                      
                      <option value="/docs/copilot/concepts/language-models" >Language Models</option>
                      
                      <option value="/docs/copilot/concepts/context" >Context</option>
                      
                      <option value="/docs/copilot/concepts/tools" >Tools</option>
                      
                      <option value="/docs/copilot/concepts/agents" >Agents</option>
                      
                      <option value="/docs/copilot/concepts/customization" >Customization</option>
                      
                      <option value="/docs/copilot/concepts/trust-and-safety" >Trust & Safety</option>
                      
                    </optgroup>
            
                    <optgroup label="GitHub Copilot - Agents">
            
                      <option value="/docs/copilot/agents/overview" >Overview</option>
                      
                      <option value="/docs/copilot/agents/agents-tutorial" >Agents Tutorial</option>
                      
                      <option value="/docs/copilot/agents/planning" >Planning</option>
                      
                      <option value="/docs/copilot/agents/memory" >Memory</option>
                      
                      <option value="/docs/copilot/agents/agent-tools" >Tools</option>
                      
                      <option value="/docs/copilot/agents/subagents" >Subagents</option>
                      
                      <option value="/docs/copilot/agents/local-agents" >Local Agents</option>
                      
                      <option value="/docs/copilot/agents/copilot-cli" >Copilot CLI</option>
                      
                      <option value="/docs/copilot/agents/cloud-agents" >Cloud Agents</option>
                      
                      <option value="/docs/copilot/agents/third-party-agents" >Third-Party Agents</option>
                      
                    </optgroup>
            
                    <optgroup label="GitHub Copilot - Chat">
            
                      <option value="/docs/copilot/chat/copilot-chat" >Overview</option>
                      
                      <option value="/docs/copilot/chat/chat-sessions" >Chat Sessions</option>
                      
                      <option value="/docs/copilot/chat/copilot-chat-context" >Add Context</option>
                      
                      <option value="/docs/copilot/chat/inline-chat" >Inline Chat</option>
                      
                      <option value="/docs/copilot/chat/review-code-edits" >Review Edits</option>
                      
                      <option value="/docs/copilot/chat/chat-checkpoints" >Checkpoints</option>
                      
                      <option value="/docs/copilot/chat/chat-artifacts" >Artifacts Panel</option>
                      
                      <option value="/docs/copilot/chat/chat-debug-view" >Debug Chat Interactions</option>
                      
                      <option value="/docs/copilot/chat/prompt-examples" >Prompt Examples</option>
                      
                    </optgroup>
            
                    <optgroup label="GitHub Copilot - Customization">
            
                      <option value="/docs/copilot/customization/overview" >Overview</option>
                      
                      <option value="/docs/copilot/customization/custom-instructions" >Instructions</option>
                      
                      <option value="/docs/copilot/customization/prompt-files" >Prompt Files</option>
                      
                      <option value="/docs/copilot/customization/custom-agents" >Custom Agents</option>
                      
                      <option value="/docs/copilot/customization/agent-skills" selected>Agent Skills</option>
                      
                      <option value="/docs/copilot/customization/language-models" >Language Models</option>
                      
                      <option value="/docs/copilot/customization/mcp-servers" >MCP</option>
                      
                      <option value="/docs/copilot/customization/hooks" >Hooks</option>
                      
                      <option value="/docs/copilot/customization/agent-plugins" >Plugins</option>
                      
                    </optgroup>
            
                    <optgroup label="GitHub Copilot - Guides & Tutorials">
            
                      <option value="/docs/copilot/guides/context-engineering-guide" >Context Engineering</option>
                      
                      <option value="/docs/copilot/guides/customize-copilot-guide" >Customize AI</option>
                      
                      <option value="/docs/copilot/guides/test-driven-development-guide" >Test-Driven Development</option>
                      
                      <option value="/docs/copilot/guides/notebooks-with-ai" >Edit Notebooks with AI</option>
                      
                      <option value="/docs/copilot/guides/test-with-copilot" >Test with AI</option>
                      
                      <option value="/docs/copilot/guides/browser-agent-testing-guide" >Test Web Apps with Browser Tools</option>
                      
                      <option value="/docs/copilot/guides/debug-with-copilot" >Debug with AI</option>
                      
                      <option value="/docs/copilot/guides/mcp-developer-guide" >MCP Dev Guide</option>
                      
                      <option value="/docs/copilot/guides/monitoring-agents" >Monitoring</option>
                      
                    </optgroup>
            
                    <option value="/docs/copilot/ai-powered-suggestions" >Inline Suggestions</option>
                    
                    <option value="/docs/copilot/copilot-smart-actions" >Smart Actions</option>
                    
                    <option value="/docs/copilot/best-practices" >Best Practices</option>
                    
                    <option value="/docs/copilot/security" >Security</option>
                    
                    <option value="/docs/copilot/troubleshooting" >Troubleshooting</option>
                    
                    <option value="/docs/copilot/faq" >FAQ</option>
                    
                    <optgroup label="GitHub Copilot - Reference">
            
                      <option value="/docs/copilot/reference/copilot-vscode-features" >Cheat Sheet</option>
                      
                      <option value="/docs/copilot/reference/copilot-settings" >Settings Reference</option>
                      
                      <option value="/docs/copilot/reference/mcp-configuration" >MCP Configuration</option>
                      
                      <option value="/docs/copilot/reference/workspace-context" >Workspace Context</option>
                      
                    </optgroup>
            
              </optgroup>
              
              <optgroup label="Configure">
            
                    <option value="/docs/configure/locales" >Display Language</option>
                    
                    <option value="/docs/configure/custom-layout" >Layout</option>
                    
                    <option value="/docs/configure/keybindings" >Keyboard Shortcuts</option>
                    
                    <option value="/docs/configure/settings" >Settings</option>
                    
                    <option value="/docs/configure/settings-sync" >Settings Sync</option>
                    
                    <optgroup label="Configure - Extensions">
            
                      <option value="/docs/configure/extensions/extension-marketplace" >Extension Marketplace</option>
                      
                      <option value="/docs/configure/extensions/extension-runtime-security" >Extension Runtime Security</option>
                      
                    </optgroup>
            
                    <option value="/docs/configure/themes" >Themes</option>
                    
                    <option value="/docs/configure/profiles" >Profiles</option>
                    
                    <optgroup label="Configure - Accessibility">
            
                      <option value="/docs/configure/accessibility/accessibility" >Overview</option>
                      
                      <option value="/docs/configure/accessibility/voice" >Voice Interactions</option>
                      
                    </optgroup>
            
                    <option value="/docs/configure/command-line" >Command Line Interface</option>
                    
                    <option value="/docs/configure/telemetry" >Telemetry</option>
                    
              </optgroup>
              
              <optgroup label="Edit Code">
            
                    <option value="/docs/editing/codebasics" >Basic Editing</option>
                    
                    <option value="/docs/editing/intellisense" >IntelliSense</option>
                    
                    <option value="/docs/editing/editingevolved" >Code Navigation</option>
                    
                    <option value="/docs/editing/refactoring" >Refactoring</option>
                    
                    <option value="/docs/editing/userdefinedsnippets" >Snippets</option>
                    
                    <optgroup label="Edit Code - Workspaces">
            
                      <option value="/docs/editing/workspaces/workspaces" >Overview</option>
                      
                      <option value="/docs/editing/workspaces/multi-root-workspaces" >Multi-Root Workspaces</option>
                      
                      <option value="/docs/editing/workspaces/workspace-trust" >Workspace Trust</option>
                      
                    </optgroup>
            
              </optgroup>
              
              <optgroup label="Build, Debug, Test">
            
                    <option value="/docs/debugtest/tasks" >Tasks</option>
                    
                    <option value="/docs/debugtest/debugging" >Debugging</option>
                    
                    <option value="/docs/debugtest/debugging-configuration" >Debug Configuration</option>
                    
                    <option value="/docs/debugtest/testing" >Testing</option>
                    
                    <option value="/docs/debugtest/port-forwarding" >Port Forwarding</option>
                    
                    <option value="/docs/debugtest/integrated-browser" >Integrated Browser</option>
                    
              </optgroup>
              
              <optgroup label="Source Control">
            
                    <option value="/docs/sourcecontrol/overview" >Overview</option>
                    
                    <option value="/docs/sourcecontrol/quickstart" >Quickstart</option>
                    
                    <option value="/docs/sourcecontrol/staging-commits" >Staging & Committing</option>
                    
                    <option value="/docs/sourcecontrol/branches-worktrees" >Branches & Worktrees</option>
                    
                    <option value="/docs/sourcecontrol/repos-remotes" >Repositories & Remotes</option>
                    
                    <option value="/docs/sourcecontrol/merge-conflicts" >Merge Conflicts</option>
                    
                    <option value="/docs/sourcecontrol/github" >Collaborate on GitHub</option>
                    
                    <option value="/docs/sourcecontrol/troubleshooting" >Troubleshooting</option>
                    
                    <option value="/docs/sourcecontrol/faq" >FAQ</option>
                    
              </optgroup>
              
              <optgroup label="Terminal">
            
                    <option value="/docs/terminal/getting-started" >Getting Started Tutorial</option>
                    
                    <option value="/docs/terminal/basics" >Terminal Basics</option>
                    
                    <option value="/docs/terminal/profiles" >Terminal Profiles</option>
                    
                    <option value="/docs/terminal/shell-integration" >Shell Integration</option>
                    
                    <option value="/docs/terminal/appearance" >Appearance</option>
                    
                    <option value="/docs/terminal/advanced" >Advanced</option>
                    
              </optgroup>
              
              <optgroup label="Enterprise">
            
                    <option value="/docs/enterprise/overview" >Overview</option>
                    
                    <option value="/docs/enterprise/policies" >Enterprise Policies</option>
                    
                    <option value="/docs/enterprise/ai-settings" >AI Settings</option>
                    
                    <option value="/docs/enterprise/extensions" >Extensions</option>
                    
                    <option value="/docs/enterprise/telemetry" >Telemetry</option>
                    
                    <option value="/docs/enterprise/updates" >Updates</option>
                    
              </optgroup>
              
              <optgroup label="Languages">
            
                    <option value="/docs/languages/overview" >Overview</option>
                    
                    <option value="/docs/languages/javascript" >JavaScript</option>
                    
                    <option value="/docs/languages/json" >JSON</option>
                    
                    <option value="/docs/languages/html" >HTML</option>
                    
                    <option value="/docs/languages/emmet" >Emmet</option>
                    
                    <option value="/docs/languages/css" >CSS, SCSS and Less</option>
                    
                    <option value="/docs/languages/typescript" >TypeScript</option>
                    
                    <option value="/docs/languages/markdown" >Markdown</option>
                    
                    <option value="/docs/languages/powershell" >PowerShell</option>
                    
                    <option value="/docs/languages/cpp" >C++</option>
                    
                    <option value="/docs/languages/java" >Java</option>
                    
                    <option value="/docs/languages/php" >PHP</option>
                    
                    <option value="/docs/languages/python" >Python</option>
                    
                    <option value="/docs/languages/julia" >Julia</option>
                    
                    <option value="/docs/languages/r" >R</option>
                    
                    <option value="/docs/languages/ruby" >Ruby</option>
                    
                    <option value="/docs/languages/rust" >Rust</option>
                    
                    <option value="/docs/languages/go" >Go</option>
                    
                    <option value="/docs/languages/tsql" >T-SQL</option>
                    
                    <option value="/docs/languages/csharp" >C#</option>
                    
                    <option value="/docs/languages/dotnet" >.NET</option>
                    
                    <option value="/docs/languages/swift" >Swift</option>
                    
              </optgroup>
              
              <optgroup label="Node.js / JavaScript">
            
                    <option value="/docs/nodejs/working-with-javascript" >Working with JavaScript</option>
                    
                    <option value="/docs/nodejs/nodejs-tutorial" >Node.js Tutorial</option>
                    
                    <option value="/docs/nodejs/nodejs-debugging" >Node.js Debugging</option>
                    
                    <option value="/docs/nodejs/nodejs-deployment" >Deploy Node.js Apps</option>
                    
                    <option value="/docs/nodejs/browser-debugging" >Browser Debugging</option>
                    
                    <option value="/docs/nodejs/angular-tutorial" >Angular Tutorial</option>
                    
                    <option value="/docs/nodejs/reactjs-tutorial" >React Tutorial</option>
                    
                    <option value="/docs/nodejs/vuejs-tutorial" >Vue Tutorial</option>
                    
                    <option value="/docs/nodejs/debugging-recipes" >Debugging Recipes</option>
                    
                    <option value="/docs/nodejs/profiling" >Performance Profiling</option>
                    
                    <option value="/docs/nodejs/extensions" >Extensions</option>
                    
              </optgroup>
              
              <optgroup label="TypeScript">
            
                    <option value="/docs/typescript/typescript-tutorial" >Tutorial</option>
                    
                    <option value="/docs/typescript/typescript-transpiling" >Transpiling</option>
                    
                    <option value="/docs/typescript/typescript-editing" >Editing</option>
                    
                    <option value="/docs/typescript/typescript-refactoring" >Refactoring</option>
                    
                    <option value="/docs/typescript/typescript-debugging" >Debugging</option>
                    
              </optgroup>
              
              <optgroup label="Python">
            
                    <option value="/docs/python/python-quick-start" >Quick Start</option>
                    
                    <option value="/docs/python/python-tutorial" >Tutorial</option>
                    
                    <option value="/docs/python/run" >Run Python Code</option>
                    
                    <option value="/docs/python/editing" >Editing</option>
                    
                    <option value="/docs/python/linting" >Linting</option>
                    
                    <option value="/docs/python/formatting" >Formatting</option>
                    
                    <option value="/docs/python/debugging" >Debugging</option>
                    
                    <option value="/docs/python/environments" >Environments</option>
                    
                    <option value="/docs/python/testing" >Testing</option>
                    
                    <option value="/docs/python/jupyter-support-py" >Python Interactive</option>
                    
                    <option value="/docs/python/tutorial-django" >Django Tutorial</option>
                    
                    <option value="/docs/python/tutorial-fastapi" >FastAPI Tutorial</option>
                    
                    <option value="/docs/python/tutorial-flask" >Flask Tutorial</option>
                    
                    <option value="/docs/python/tutorial-create-containers" >Create Containers</option>
                    
                    <option value="/docs/python/python-on-azure" >Deploy Python Apps</option>
                    
                    <option value="/docs/python/python-web" >Python in the Web</option>
                    
                    <option value="/docs/python/settings-reference" >Settings Reference</option>
                    
              </optgroup>
              
              <optgroup label="Java">
            
                    <option value="/docs/java/java-tutorial" >Getting Started</option>
                    
                    <option value="/docs/java/java-editing" >Navigate and Edit</option>
                    
                    <option value="/docs/java/java-refactoring" >Refactoring</option>
                    
                    <option value="/docs/java/java-linting" >Formatting and Linting</option>
                    
                    <option value="/docs/java/java-project" >Project Management</option>
                    
                    <option value="/docs/java/java-build" >Build Tools</option>
                    
                    <option value="/docs/java/java-debugging" >Run and Debug</option>
                    
                    <option value="/docs/java/java-testing" >Testing</option>
                    
                    <option value="/docs/java/java-spring-boot" >Spring Boot</option>
                    
                    <option value="/docs/java/java-app-mod" >Modernizing Java Apps</option>
                    
                    <option value="/docs/java/java-tomcat-jetty" >Application Servers</option>
                    
                    <option value="/docs/java/java-on-azure" >Deploy Java Apps</option>
                    
                    <option value="/docs/java/java-gui" >GUI Applications</option>
                    
                    <option value="/docs/java/extensions" >Extensions</option>
                    
                    <option value="/docs/java/java-faq" >FAQ</option>
                    
              </optgroup>
              
              <optgroup label="C++">
            
                    <option value="/docs/cpp/introvideos-cpp" >Intro Videos</option>
                    
                    <option value="/docs/cpp/config-linux" >GCC on Linux</option>
                    
                    <option value="/docs/cpp/config-mingw" >GCC on Windows</option>
                    
                    <option value="/docs/cpp/config-wsl" >GCC on Windows Subsystem for Linux</option>
                    
                    <option value="/docs/cpp/config-clang-mac" >Clang on macOS</option>
                    
                    <option value="/docs/cpp/config-msvc" >Microsoft C++ on Windows</option>
                    
                    <option value="/docs/cpp/build-with-cmake" >Build with CMake</option>
                    
                    <option value="/docs/cpp/cmake-linux" >CMake Tools on Linux</option>
                    
                    <option value="/docs/cpp/cmake-quickstart" >CMake Quick Start</option>
                    
                    <option value="/docs/cpp/cpp-devtools" >C++ Dev Tools for Copilot</option>
                    
                    <option value="/docs/cpp/cpp-ide" >Editing and Navigating</option>
                    
                    <option value="/docs/cpp/cpp-debug" >Debugging</option>
                    
                    <option value="/docs/cpp/launch-json-reference" >Configure Debugging</option>
                    
                    <option value="/docs/cpp/cpp-refactoring" >Refactoring</option>
                    
                    <option value="/docs/cpp/customize-cpp-settings" >Settings Reference</option>
                    
                    <option value="/docs/cpp/configure-intellisense" >Configure IntelliSense</option>
                    
                    <option value="/docs/cpp/configure-intellisense-crosscompilation" >Configure IntelliSense for Cross-Compiling</option>
                    
                    <option value="/docs/cpp/faq-cpp" >FAQ</option>
                    
              </optgroup>
              
              <optgroup label="C#">
            
                    <option value="/docs/csharp/introvideos-csharp" >Intro Videos</option>
                    
                    <option value="/docs/csharp/get-started" >Get Started</option>
                    
                    <option value="/docs/csharp/navigate-edit" >Navigate and Edit</option>
                    
                    <option value="/docs/csharp/intellicode" >IntelliCode</option>
                    
                    <option value="/docs/csharp/refactoring" >Refactoring</option>
                    
                    <option value="/docs/csharp/formatting-linting" >Formatting and Linting</option>
                    
                    <option value="/docs/csharp/project-management" >Project Management</option>
                    
                    <option value="/docs/csharp/build-tools" >Build Tools</option>
                    
                    <option value="/docs/csharp/package-management" >Package Management</option>
                    
                    <option value="/docs/csharp/debugging" >Run and Debug</option>
                    
                    <option value="/docs/csharp/testing" >Testing</option>
                    
                    <option value="/docs/csharp/cs-dev-kit-faq" >FAQ</option>
                    
              </optgroup>
              
              <optgroup label="Container Tools">
            
                    <option value="/docs/containers/overview" >Overview</option>
                    
                    <option value="/docs/containers/quickstart-node" >Node.js</option>
                    
                    <option value="/docs/containers/quickstart-python" >Python</option>
                    
                    <option value="/docs/containers/quickstart-aspnet-core" >ASP.NET Core</option>
                    
                    <option value="/docs/containers/debug-common" >Debug</option>
                    
                    <option value="/docs/containers/docker-compose" >Docker Compose</option>
                    
                    <option value="/docs/containers/quickstart-container-registries" >Registries</option>
                    
                    <option value="/docs/containers/app-service" >Deploy to Azure</option>
                    
                    <option value="/docs/containers/choosing-dev-environment" >Choose a Dev Environment</option>
                    
                    <option value="/docs/containers/reference" >Customize</option>
                    
                    <option value="/docs/containers/bridge-to-kubernetes" >Develop with Kubernetes</option>
                    
                    <option value="/docs/containers/troubleshooting" >Tips and Tricks</option>
                    
              </optgroup>
              
              <optgroup label="Data Science">
            
                    <option value="/docs/datascience/overview" >Overview</option>
                    
                    <option value="/docs/datascience/jupyter-notebooks" >Jupyter Notebooks</option>
                    
                    <option value="/docs/datascience/data-science-tutorial" >Data Science Tutorial</option>
                    
                    <option value="/docs/datascience/python-interactive" >Python Interactive</option>
                    
                    <option value="/docs/datascience/data-wrangler-quick-start" >Data Wrangler Quick Start</option>
                    
                    <option value="/docs/datascience/data-wrangler" >Data Wrangler</option>
                    
                    <option value="/docs/datascience/pytorch-support" >PyTorch Support</option>
                    
                    <option value="/docs/datascience/azure-machine-learning" >Azure Machine Learning</option>
                    
                    <option value="/docs/datascience/jupyter-kernel-management" >Manage Jupyter Kernels</option>
                    
                    <option value="/docs/datascience/notebooks-web" >Jupyter Notebooks on the Web</option>
                    
                    <option value="/docs/datascience/microsoft-fabric-quickstart" >Data Science in Microsoft Fabric</option>
                    
              </optgroup>
              
              <optgroup label="Intelligent Apps">
            
                    <option value="/docs/intelligentapps/overview" >Foundry Toolkit Overview</option>
                    
                    <option value="/docs/intelligentapps/copilot-tools" >Foundry Toolkit Copilot Tools</option>
                    
                    <option value="/docs/intelligentapps/create-agents" >Create Agents</option>
                    
                    <option value="/docs/intelligentapps/models" >Models</option>
                    
                    <option value="/docs/intelligentapps/playground" >Playground</option>
                    
                    <option value="/docs/intelligentapps/agentbuilder" >Agent Builder</option>
                    
                    <option value="/docs/intelligentapps/agent-inspector" >Agent Inspector</option>
                    
                    <option value="/docs/intelligentapps/evaluation" >Evaluation</option>
                    
                    <option value="/docs/intelligentapps/tool-catalog" >Tool Catalog</option>
                    
                    <option value="/docs/intelligentapps/finetune" >Fine-Tuning (Automated Setup)</option>
                    
                    <option value="/docs/intelligentapps/finetune-legacy" >Fine-Tuning (Project Template)</option>
                    
                    <option value="/docs/intelligentapps/modelconversion" >Model Conversion</option>
                    
                    <option value="/docs/intelligentapps/tracing" >Tracing</option>
                    
                    <option value="/docs/intelligentapps/profiling" >Profiling (Windows ML)</option>
                    
                    <option value="/docs/intelligentapps/faq" >FAQ</option>
                    
                    <optgroup label="Intelligent Apps - Reference">
            
                      <option value="/docs/intelligentapps/reference/FileStructure" >File Structure</option>
                      
                      <option value="/docs/intelligentapps/reference/ManualModelConversion" >Manual Model Conversion</option>
                      
                      <option value="/docs/intelligentapps/reference/ManualConversionOnGPU" >Manual Model Conversion on GPU</option>
                      
                      <option value="/docs/intelligentapps/reference/SetupWithoutAITK" >Setup Environment Without Foundry Toolkit</option>
                      
                      <option value="/docs/intelligentapps/reference/TemplateProject" >Template Project</option>
                      
                      <option value="/docs/intelligentapps/reference/migrate-from-visualizer" >Migrating from Visualizer to Agent Inspector</option>
                      
                    </optgroup>
            
              </optgroup>
              
              <optgroup label="Azure">
            
                    <option value="/docs/azure/overview" >Overview</option>
                    
                    <option value="/docs/azure/gettingstarted" >Getting Started</option>
                    
                    <option value="/docs/azure/resourcesextension" >Resources View</option>
                    
                    <option value="/docs/azure/deployment" >Deployment</option>
                    
                    <option value="/docs/azure/vscodeforweb" >VS Code for the Web - Azure</option>
                    
                    <option value="/docs/azure/containers" >Containers</option>
                    
                    <option value="/docs/azure/aksextensions" >Azure Kubernetes Service</option>
                    
                    <option value="/docs/azure/kubernetes" >Kubernetes</option>
                    
                    <option value="/docs/azure/mongodb" >MongoDB</option>
                    
                    <option value="/docs/azure/remote-debugging" >Remote Debugging for Node.js</option>
                    
              </optgroup>
              
              <optgroup label="Remote">
            
                    <option value="/docs/remote/remote-overview" >Overview</option>
                    
                    <option value="/docs/remote/ssh" >SSH</option>
                    
                    <option value="/docs/remote/dev-containers" >Dev Containers</option>
                    
                    <option value="/docs/remote/wsl" >Windows Subsystem for Linux</option>
                    
                    <option value="/docs/remote/codespaces" >GitHub Codespaces</option>
                    
                    <option value="/docs/remote/vscode-server" >VS Code Server</option>
                    
                    <option value="/docs/remote/tunnels" >Tunnels</option>
                    
                    <option value="/docs/remote/ssh-tutorial" >SSH Tutorial</option>
                    
                    <option value="/docs/remote/wsl-tutorial" >WSL Tutorial</option>
                    
                    <option value="/docs/remote/troubleshooting" >Tips and Tricks</option>
                    
                    <option value="/docs/remote/faq" >FAQ</option>
                    
              </optgroup>
              
              <optgroup label="Dev Containers">
            
                    <option value="/docs/devcontainers/containers" >Overview</option>
                    
                    <option value="/docs/devcontainers/tutorial" >Tutorial</option>
                    
                    <option value="/docs/devcontainers/attach-container" >Attach to Container</option>
                    
                    <option value="/docs/devcontainers/create-dev-container" >Create Dev Container</option>
                    
                    <option value="/docs/devcontainers/containers-advanced" >Advanced Containers</option>
                    
                    <option value="/docs/devcontainers/devcontainerjson-reference" >devcontainer.json</option>
                    
                    <option value="/docs/devcontainers/devcontainer-cli" >Dev Container CLI</option>
                    
                    <option value="/docs/devcontainers/tips-and-tricks" >Tips and Tricks</option>
                    
                    <option value="/docs/devcontainers/faq" >FAQ</option>
                    
              </optgroup>
              
              <optgroup label="Reference">
            
                    <option value="/docs/reference/default-keybindings" >Default Keyboard Shortcuts</option>
                    
                    <option value="/docs/reference/default-settings" >Default Settings</option>
                    
                    <option value="/docs/reference/variables-reference" >Substitution Variables</option>
                    
                    <option value="/docs/reference/tasks-appendix" >Tasks Schema</option>
                    
              </optgroup>
              
              </select>
            </nav>        </aside>
        
        <!-- Content wrapper contains main content + right sidebar -->
        <div class="docs-content-wrapper">
            <!-- Right sidebar - On This Page -->
            <aside class="docs-right-sidebar hidden-xs">
                <div class="docs-markdown-actions">
                    <div class="docs-markdown-dropdown" data-raw-url="/raw/docs/copilot/customization/agent-skills.md">
                        <button type="button" class="docs-markdown-btn-main" data-action="copy-markdown" aria-label="Copy as Markdown">
                            <span class="codicon codicon-copy" aria-hidden="true"></span>
                            <span>Copy as Markdown</span>
                        </button>
                        <button type="button" class="docs-markdown-btn-trigger" aria-haspopup="true" aria-expanded="false" aria-label="More Markdown options">
                            <span class="codicon codicon-chevron-down docs-markdown-chevron" aria-hidden="true"></span>
                        </button>
                        <ul class="docs-markdown-menu" role="menu" aria-label="Markdown options">
                            <li role="menuitem" tabindex="0" data-action="copy-markdown">
                                <span class="codicon codicon-copy" aria-hidden="true"></span>
                                <span>Copy as Markdown</span>
                            </li>
                            <li role="menuitem" tabindex="0" data-action="view-markdown">
                                <span class="codicon codicon-file" aria-hidden="true"></span>
                                <span>View as Markdown</span>
                                <span class="codicon codicon-link-external" aria-hidden="true"></span>
                            </li>
                        </ul>
                    </div>
                </div>
                <nav id="docs-subnavbar" aria-label="On Page">
                    
                    <h4><span class="sr-only">On this page there are 10 sections</span><span
                            aria-hidden="true">On this page</span></h4>
                    <ul class="nav">
                        
                        <li><a href="#_agent-skills-vs-custom-instructions">Agent Skills vs custom instructions</a></li>
                        
                        <li><a href="#_create-a-skill">Create a skill</a></li>
                        
                        <li><a href="#_skillmd-file-format">SKILL.md file format</a></li>
                        
                        <li><a href="#_example-skills">Example skills</a></li>
                        
                        <li><a href="#_use-skills-as-slash-commands">Use skills as slash commands</a></li>
                        
                        <li><a href="#_how-copilot-uses-skills">How Copilot uses skills</a></li>
                        
                        <li><a href="#_use-shared-skills">Use shared skills</a></li>
                        
                        <li><a href="#_contribute-skills-from-extensions">Contribute skills from extensions</a></li>
                        
                        <li><a href="#_agent-skills-standard">Agent Skills standard</a></li>
                        
                        <li><a href="#_related-resources">Related resources</a></li>
                        
                    </ul>
                    
                </nav>
                
            </aside>
            
            <!-- Main article content -->
            <main class="docs-main-content body">
                <h1>Use Agent Skills in VS Code</h1>
<p>Agent Skills are folders of instructions, scripts, and resources that GitHub Copilot can load when relevant to perform specialized tasks. Agent Skills is an <a href="https://agentskills.io" class="external-link" target="_blank">open standard</a> that works across multiple AI agents, including GitHub Copilot in VS Code, GitHub Copilot CLI, and GitHub Copilot cloud agent.</p>
<p>Unlike <a href="/docs/copilot/customization/custom-instructions">custom instructions</a> that primarily define coding guidelines, skills enable specialized capabilities and workflows that can include scripts, examples, and other resources. Skills you create are portable and work across any skills-compatible agent.</p>
<p>Key benefits of Agent Skills:</p>
<ul>
<li><strong>Specialize Copilot</strong>: Tailor capabilities for domain-specific tasks without repeating context</li>
<li><strong>Reduce repetition</strong>: Create once, use automatically across all conversations</li>
<li><strong>Compose capabilities</strong>: Combine multiple skills to build complex workflows</li>
<li><strong>Efficient loading</strong>: Only relevant content loads into context when needed</li>
</ul>
<div class="markdown-alert tip" dir="auto">
      <span>
        <svg class="markdown-alert-icon" viewBox="0 0 16 16" version="1.1" width="16" height="16" aria-hidden="true">
          <path d="M8 1.5c-2.363 0-4 1.69-4 3.75 0 .984.424 1.625.984 2.304l.214.253c.223.264.47.556.673.848.284.411.537.896.621 1.49a.75.75 0 0 1-1.484.211c-.04-.282-.163-.547-.37-.847a8.456 8.456 0 0 0-.542-.68c-.084-.1-.173-.205-.268-.32C3.201 7.75 2.5 6.766 2.5 5.25 2.5 2.31 4.863 0 8 0s5.5 2.31 5.5 5.25c0 1.516-.701 2.5-1.328 3.259-.095.115-.184.22-.268.319-.207.245-.383.453-.541.681-.208.3-.33.565-.37.847a.751.751 0 0 1-1.485-.212c.084-.593.337-1.078.621-1.489.203-.292.45-.584.673-.848.075-.088.147-.173.213-.253.561-.679.985-1.32.985-2.304 0-2.06-1.637-3.75-4-3.75ZM5.75 12h4.5a.75.75 0 0 1 0 1.5h-4.5a.75.75 0 0 1 0-1.5ZM6 15.25a.75.75 0 0 1 .75-.75h2.5a.75.75 0 0 1 0 1.5h-2.5a.75.75 0 0 1-.75-.75Z">
          </path>
        </svg>
        Tip
      </span><p>Use the <a href="/docs/copilot/customization/overview#_chat-customizations-editor">Chat Customizations editor</a> (Preview) to discover, create, and manage all your chat customizations in one place. Run <strong>Chat: Open Chat Customizations</strong> from the Command Palette.</p>
</div><h2 id="_agent-skills-vs-custom-instructions" data-needslink="_agent-skills-vs-custom-instructions">Agent Skills vs custom instructions</h2>
<p>While both Agent Skills and custom instructions help customize Copilot's behavior, they serve different purposes:</p>
<table class="table table-striped">
<thead>
<tr>
<th>Feature</th>
<th>Agent Skills</th>
<th>Custom Instructions</th>
</tr>
</thead>
<tbody>
<tr>
<td><strong>Purpose</strong></td>
<td>Teach specialized capabilities and workflows</td>
<td>Define coding standards and guidelines</td>
</tr>
<tr>
<td><strong>Portability</strong></td>
<td>Works across VS Code, Copilot CLI, and Copilot cloud agent</td>
<td>VS Code and GitHub.com only</td>
</tr>
<tr>
<td><strong>Content</strong></td>
<td>Instructions, scripts, examples, and resources</td>
<td>Instructions only</td>
</tr>
<tr>
<td><strong>Scope</strong></td>
<td>Task-specific, loaded on-demand</td>
<td>Always applied (or via glob patterns)</td>
</tr>
<tr>
<td><strong>Standard</strong></td>
<td>Open standard (<a href="https://agentskills.io" class="external-link" target="_blank">agentskills.io</a>)</td>
<td>VS Code-specific</td>
</tr>
</tbody>
</table>
<p>Use Agent Skills when you want to:</p>
<ul>
<li>Create reusable capabilities that work across different AI tools</li>
<li>Include scripts, examples, or other resources alongside instructions</li>
<li>Share capabilities with the wider AI community</li>
<li>Define specialized workflows like testing, debugging, or deployment processes</li>
</ul>
<p>Use custom instructions when you want to:</p>
<ul>
<li>Define project-specific coding standards</li>
<li>Set language or framework conventions</li>
<li>Specify code review or commit message guidelines</li>
<li>Apply rules based on file types using glob patterns</li>
</ul>
<h2 id="_create-a-skill" data-needslink="_create-a-skill">Create a skill</h2>
<div class="markdown-alert tip" dir="auto">
      <span>
        <svg class="markdown-alert-icon" viewBox="0 0 16 16" version="1.1" width="16" height="16" aria-hidden="true">
          <path d="M8 1.5c-2.363 0-4 1.69-4 3.75 0 .984.424 1.625.984 2.304l.214.253c.223.264.47.556.673.848.284.411.537.896.621 1.49a.75.75 0 0 1-1.484.211c-.04-.282-.163-.547-.37-.847a8.456 8.456 0 0 0-.542-.68c-.084-.1-.173-.205-.268-.32C3.201 7.75 2.5 6.766 2.5 5.25 2.5 2.31 4.863 0 8 0s5.5 2.31 5.5 5.25c0 1.516-.701 2.5-1.328 3.259-.095.115-.184.22-.268.319-.207.245-.383.453-.541.681-.208.3-.33.565-.37.847a.751.751 0 0 1-1.485-.212c.084-.593.337-1.078.621-1.489.203-.292.45-.584.673-.848.075-.088.147-.173.213-.253.561-.679.985-1.32.985-2.304 0-2.06-1.637-3.75-4-3.75ZM5.75 12h4.5a.75.75 0 0 1 0 1.5h-4.5a.75.75 0 0 1 0-1.5ZM6 15.25a.75.75 0 0 1 .75-.75h2.5a.75.75 0 0 1 0 1.5h-2.5a.75.75 0 0 1-.75-.75Z">
          </path>
        </svg>
        Tip
      </span><p>Type <code>/skills</code> in the chat input to quickly open the <strong>Configure Skills</strong> menu.</p>
</div><p>Skills are stored in directories with a <code>SKILL.md</code> file that defines the skill's behavior. VS Code supports two types of skills:</p>
<table class="table table-striped">
<thead>
<tr>
<th>Skill type</th>
<th>Location</th>
</tr>
</thead>
<tbody>
<tr>
<td>Project skills, stored in your repository</td>
<td><code>.github/skills/</code>, <code>.claude/skills/</code>, <code>.agents/skills/</code></td>
</tr>
<tr>
<td>Personal skills, stored in your user profile</td>
<td><code>~/.copilot/skills/</code>, <code>~/.claude/skills/</code>, <code>~/.agents/skills/</code></td>
</tr>
</tbody>
</table>
<p>You can configure additional file locations for project skills with the <span class="setting"><span class="setting-dropdown" data-setting-id="chat.agentSkillsLocations">
    <span class="setting-link-main">
      <span class="codicon codicon-settings-gear dynamic-setting-icon"></span>
      chat.agentSkillsLocations
    </span>
    <button 
      type="button"
      class="setting-dropdown-trigger"
      aria-haspopup="menu"
      aria-expanded="false"
      
      aria-label="Open setting in VS Code">
        <svg class="setting-dropdown-chevron" viewBox="0 0 16 16" fill="currentColor" aria-hidden="true">
        <path d="M7.976 10.072l4.357-4.357.62.618L8.284 11h-.618L3 6.333l.619-.618 4.357 4.357z"/>
        </svg>
    </button>
    <span class="setting-dropdown-menu" role="menu" aria-label="Open setting">
      <span role="menuitem" class="setting-action-item" data-version="stable" tabindex="0">Open in VS Code</span>
      <span role="menuitem" class="setting-action-item" data-version="insiders" tabindex="-1">Open in VS Code Insiders</span>
    </span>
  </span></span> setting. This is useful if you want to organize skills in a different folder structure or have multiple skill directories.</p>
<div class="markdown-alert tip" dir="auto">
      <span>
        <svg class="markdown-alert-icon" viewBox="0 0 16 16" version="1.1" width="16" height="16" aria-hidden="true">
          <path d="M8 1.5c-2.363 0-4 1.69-4 3.75 0 .984.424 1.625.984 2.304l.214.253c.223.264.47.556.673.848.284.411.537.896.621 1.49a.75.75 0 0 1-1.484.211c-.04-.282-.163-.547-.37-.847a8.456 8.456 0 0 0-.542-.68c-.084-.1-.173-.205-.268-.32C3.201 7.75 2.5 6.766 2.5 5.25 2.5 2.31 4.863 0 8 0s5.5 2.31 5.5 5.25c0 1.516-.701 2.5-1.328 3.259-.095.115-.184.22-.268.319-.207.245-.383.453-.541.681-.208.3-.33.565-.37.847a.751.751 0 0 1-1.485-.212c.084-.593.337-1.078.621-1.489.203-.292.45-.584.673-.848.075-.088.147-.173.213-.253.561-.679.985-1.32.985-2.304 0-2.06-1.637-3.75-4-3.75ZM5.75 12h4.5a.75.75 0 0 1 0 1.5h-4.5a.75.75 0 0 1 0-1.5ZM6 15.25a.75.75 0 0 1 .75-.75h2.5a.75.75 0 0 1 0 1.5h-2.5a.75.75 0 0 1-.75-.75Z">
          </path>
        </svg>
        Tip
      </span><p>In a monorepo, enable <span class="setting"><span class="setting-dropdown" data-setting-id="chat.useCustomizationsInParentRepositories">
    <span class="setting-link-main">
      <span class="codicon codicon-settings-gear dynamic-setting-icon"></span>
      chat.useCustomizationsInParentRepositories
    </span>
    <button 
      type="button"
      class="setting-dropdown-trigger"
      aria-haspopup="menu"
      aria-expanded="false"
      
      aria-label="Open setting in VS Code">
        <svg class="setting-dropdown-chevron" viewBox="0 0 16 16" fill="currentColor" aria-hidden="true">
        <path d="M7.976 10.072l4.357-4.357.62.618L8.284 11h-.618L3 6.333l.619-.618 4.357 4.357z"/>
        </svg>
    </button>
    <span class="setting-dropdown-menu" role="menu" aria-label="Open setting">
      <span role="menuitem" class="setting-action-item" data-version="stable" tabindex="0">Open in VS Code</span>
      <span role="menuitem" class="setting-action-item" data-version="insiders" tabindex="-1">Open in VS Code Insiders</span>
    </span>
  </span></span> to discover skills from the parent repository root. Learn more about <a href="/docs/copilot/customization/overview#_parent-repository-discovery">parent repository discovery</a>.</p>
</div><p>To create a skill:</p>
<ol>
<li>
<p>In the Chat view, select <strong>Configure Chat</strong> (gear icon) to open the Chat Customizations editor and then select the <strong>Skills</strong> tab.</p>
</li>
<li>
<p>Select <strong>New Skill (Workspace)</strong> or <strong>New Skill (User)</strong> from the dropdown, depending on where you want to store the skill.</p>
<p><img src="/assets/docs/copilot/customization/create-skill.png" alt="Screenshot of the Chat Customizations editor, showing the Skills tab and the dropdown to create a new skill." loading="lazy"></p>
</li>
<li>
<p>Select the location and enter a name for the skill.</p>
</li>
<li>
<p>Complete the <code>SKILL.md</code> file by filling in the YAML frontmatter and adding instructions in the body of the file.</p>
<pre class="shiki" data-lang="markdown" shiki-themes dark-plus light-plus" style="--shiki-dark:#D4D4D4;--shiki-light:#000000;--shiki-dark-bg:#1E1E1E;--shiki-light-bg:#FFFFFF" tabindex="0"><code><span class="line"><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000">---</span></span>
<span class="line"><span style="--shiki-dark:#569CD6;--shiki-light:#800000">name</span><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000FF">: </span><span style="--shiki-dark:#CE9178;--shiki-light:#0000FF">skill-name</span></span>
<span class="line"><span style="--shiki-dark:#569CD6;--shiki-light:#800000">description</span><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000FF">: </span><span style="--shiki-dark:#CE9178;--shiki-light:#0000FF">Description of what the skill does and when to use it</span></span>
<span class="line"><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000">---</span></span>
<span class="line"></span>
<span class="line"><span style="--shiki-dark:#569CD6;--shiki-dark-font-weight:bold;--shiki-light:#800000;--shiki-light-font-weight:bold"># Skill Instructions</span></span>
<span class="line"></span>
<span class="line"><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000">Your detailed instructions, guidelines, and examples go here...</span></span>
<span class="line"></span></code></pre>
</li>
<li>
<p>Optionally, add scripts, examples, or other resources to your skill's directory.</p>
<p>For example, a skill for testing web applications might include:</p>
<ul>
<li><code>SKILL.md</code> - Instructions for running tests</li>
<li><code>test-template.js</code> - A template test file</li>
<li><code>examples/</code> - Example test scenarios</li>
</ul>
<div class="markdown-alert note" dir="auto">
      <span>
        <svg class="markdown-alert-icon" viewBox="0 0 16 16" version="1.1" width="16" height="16" aria-hidden="true">
          <path d="M0 8a8 8 0 1 1 16 0A8 8 0 0 1 0 8Zm8-6.5a6.5 6.5 0 1 0 0 13 6.5 6.5 0 0 0 0-13ZM6.5 7.75A.75.75 0 0 1 7.25 7h1a.75.75 0 0 1 .75.75v2.75h.25a.75.75 0 0 1 0 1.5h-2a.75.75 0 0 1 0-1.5h.25v-2h-.25a.75.75 0 0 1-.75-.75ZM8 6a1 1 0 1 1 0-2 1 1 0 0 1 0 2Z">
          </path>
        </svg>
        Note
      </span><p>Make sure to reference any additional files in your <code>SKILL.md</code> for them to be picked up by the agent. Use Markdown link syntax with relative paths, such as <code>[test template](./test-template.js)</code>.</p>
</div></li>
</ol>
<h3 id="_generate-a-skill-with-ai" data-needslink="_generate-a-skill-with-ai">Generate a skill with AI</h3>
<p>You can use AI to generate a skill based on a description of the capability. Type <code>/create-skill</code> in chat and describe the skill you want (for example, &quot;a skill for running and debugging integration tests&quot;). The agent asks clarifying questions and generates a <code>SKILL.md</code> file with the directory structure, instructions, and frontmatter.</p>
<p>You can also extract a reusable skill from an ongoing conversation. For example, after a multi-turn session where you debugged a complex issue, ask &quot;create a skill from how we just debugged that&quot; to capture the multi-step procedure as a reusable skill.</p>
<p>You can also generate a skill from the Chat Customizations editor by selecting <strong>Generate Skill</strong> from the dropdown.</p>
<h2 id="_skillmd-file-format" data-needslink="_skillmd-file-format">SKILL.md file format</h2>
<p>The <code>SKILL.md</code> file is a Markdown file with YAML frontmatter that defines the skill's metadata and behavior.</p>
<h3 id="_header-required" data-needslink="_header-required">Header (required)</h3>
<p>The header is formatted as YAML frontmatter with the following fields:</p>
<table class="table table-striped">
<thead>
<tr>
<th>Field</th>
<th>Required</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td><code>name</code></td>
<td>Yes</td>
<td>A unique identifier for the skill. Only lowercase letters, numbers, and hyphens are allowed (for example, <code>webapp-testing</code>). Do not use slashes, colons, dots, or namespace prefixes. Must match the parent directory name. Maximum 64 characters. Names with invalid characters cause the skill to silently fail to load.</td>
</tr>
<tr>
<td><code>description</code></td>
<td>Yes</td>
<td>A description of what the skill does <strong>and when to use it</strong>. Be specific about both capabilities and use cases to help Copilot decide when to load the skill. Maximum 1024 characters.</td>
</tr>
<tr>
<td><code>argument-hint</code></td>
<td>No</td>
<td>Hint text shown in the chat input field when the skill is invoked as a slash command. Helps users understand what additional information to provide (for example, <code>[test file] [options]</code>).</td>
</tr>
<tr>
<td><code>user-invocable</code></td>
<td>No</td>
<td>Controls whether the skill appears as a slash command in the chat menu. Defaults to <code>true</code>. Set to <code>false</code> to hide the skill from the <code>/</code> menu while still allowing the agent to load it automatically.</td>
</tr>
<tr>
<td><code>disable-model-invocation</code></td>
<td>No</td>
<td>Controls whether the agent can automatically load the skill based on relevance. Defaults to <code>false</code>. Set to <code>true</code> to require manual invocation through the <code>/</code> slash command only.</td>
</tr>
</tbody>
</table>
<div class="markdown-alert important" dir="auto">
      <span>
        <svg class="markdown-alert-icon" viewBox="0 0 16 16" version="1.1" width="16" height="16" aria-hidden="true">
          <path d="M0 1.75C0 .784.784 0 1.75 0h12.5C15.216 0 16 .784 16 1.75v9.5A1.75 1.75 0 0 1 14.25 13H8.06l-2.573 2.573A1.458 1.458 0 0 1 3 14.543V13H1.75A1.75 1.75 0 0 1 0 11.25Zm1.75-.25a.25.25 0 0 0-.25.25v9.5c0 .138.112.25.25.25h2a.75.75 0 0 1 .75.75v2.19l2.72-2.72a.749.749 0 0 1 .53-.22h6.5a.25.25 0 0 0 .25-.25v-9.5a.25.25 0 0 0-.25-.25Zm7 2.25v2.5a.75.75 0 0 1-1.5 0v-2.5a.75.75 0 0 1 1.5 0ZM9 9a1 1 0 1 1-2 0 1 1 0 0 1 2 0Z">
          </path>
        </svg>
        Important
      </span><p>When a skill is distributed through a <a href="/docs/copilot/customization/agent-plugins">plugin</a>, the plugin name is automatically used as a command prefix (for example, <code>/my-plugin:test-runner</code>). Do not manually add namespace prefixes to the skill <code>name</code> field. Using prefixes like <code>myorg/skillname</code> or <code>myorg:skillname</code> causes the skill to silently fail to load.</p>
</div><h3 id="_body" data-needslink="_body">Body</h3>
<p>The skill body contains the instructions, guidelines, and examples that Copilot should follow when using this skill. Write clear, specific instructions that describe:</p>
<ul>
<li>What the skill helps accomplish</li>
<li>When to use the skill</li>
<li>Step-by-step procedures to follow</li>
<li>Examples of the expected input and output</li>
<li>References to any included scripts or resources</li>
</ul>
<p>You can reference files within the skill directory using relative paths. For example, to reference a script in your skill directory, use <code>[test script](./test-template.js)</code>.</p>
<h2 id="_example-skills" data-needslink="_example-skills">Example skills</h2>
<p>The following examples demonstrate different types of skills you can create.</p>
<details>
<summary>Example: Web application testing skill</summary>
<pre class="shiki" data-lang="markdown" shiki-themes dark-plus light-plus" style="--shiki-dark:#D4D4D4;--shiki-light:#000000;--shiki-dark-bg:#1E1E1E;--shiki-light-bg:#FFFFFF" tabindex="0"><code><span class="line"><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000">---</span></span>
<span class="line"><span style="--shiki-dark:#569CD6;--shiki-light:#800000">name</span><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000FF">: </span><span style="--shiki-dark:#CE9178;--shiki-light:#0000FF">webapp-testing</span></span>
<span class="line"><span style="--shiki-dark:#569CD6;--shiki-light:#800000">description</span><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000FF">: </span><span style="--shiki-dark:#CE9178;--shiki-light:#0000FF">Guide for testing web applications using Playwright. Use this when asked to create or run browser-based tests.</span></span>
<span class="line"><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000">---</span></span>
<span class="line"></span>
<span class="line"><span style="--shiki-dark:#569CD6;--shiki-dark-font-weight:bold;--shiki-light:#800000;--shiki-light-font-weight:bold"># Web Application Testing with Playwright</span></span>
<span class="line"></span>
<span class="line"><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000">This skill helps you create and run browser-based tests for web applications using Playwright.</span></span>
<span class="line"></span>
<span class="line"><span style="--shiki-dark:#569CD6;--shiki-dark-font-weight:bold;--shiki-light:#800000;--shiki-light-font-weight:bold">## When to use this skill</span></span>
<span class="line"></span>
<span class="line"><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000">Use this skill when you need to:</span></span>
<span class="line"><span style="--shiki-dark:#6796E6;--shiki-light:#0451A5">-</span><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000"> Create new Playwright tests for web applications</span></span>
<span class="line"><span style="--shiki-dark:#6796E6;--shiki-light:#0451A5">-</span><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000"> Debug failing browser tests</span></span>
<span class="line"><span style="--shiki-dark:#6796E6;--shiki-light:#0451A5">-</span><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000"> Set up test infrastructure for a new project</span></span>
<span class="line"></span>
<span class="line"><span style="--shiki-dark:#569CD6;--shiki-dark-font-weight:bold;--shiki-light:#800000;--shiki-light-font-weight:bold">## Creating tests</span></span>
<span class="line"></span>
<span class="line"><span style="--shiki-dark:#6796E6;--shiki-light:#0451A5">1.</span><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000"> Review the [</span><span style="--shiki-dark:#CE9178;--shiki-light:#A31515">test template</span><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000">](</span><span style="--shiki-dark:#D4D4D4;--shiki-dark-text-decoration:underline;--shiki-light:#000000;--shiki-light-text-decoration:underline">./test-template.js</span><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000">) for the standard test structure</span></span>
<span class="line"><span style="--shiki-dark:#6796E6;--shiki-light:#0451A5">2.</span><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000"> Identify the user flow to test</span></span>
<span class="line"><span style="--shiki-dark:#6796E6;--shiki-light:#0451A5">3.</span><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000"> Create a new test file in the </span><span style="--shiki-dark:#CE9178;--shiki-light:#800000">`tests/`</span><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000"> directory</span></span>
<span class="line"><span style="--shiki-dark:#6796E6;--shiki-light:#0451A5">4.</span><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000"> Use Playwright's locators to find elements (prefer role-based selectors)</span></span>
<span class="line"><span style="--shiki-dark:#6796E6;--shiki-light:#0451A5">5.</span><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000"> Add assertions to verify expected behavior</span></span>
<span class="line"></span>
<span class="line"><span style="--shiki-dark:#569CD6;--shiki-dark-font-weight:bold;--shiki-light:#800000;--shiki-light-font-weight:bold">## Running tests</span></span>
<span class="line"></span>
<span class="line"><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000">To run tests locally:</span></span>
<span class="line"><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000">```bash</span></span>
<span class="line"><span style="--shiki-dark:#DCDCAA;--shiki-light:#795E26">npx</span><span style="--shiki-dark:#CE9178;--shiki-light:#A31515"> playwright</span><span style="--shiki-dark:#CE9178;--shiki-light:#A31515"> test</span></span>
<span class="line"><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000">```</span></span>
<span class="line"></span>
<span class="line"><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000">To debug tests:</span></span>
<span class="line"><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000">```bash</span></span>
<span class="line"><span style="--shiki-dark:#DCDCAA;--shiki-light:#795E26">npx</span><span style="--shiki-dark:#CE9178;--shiki-light:#A31515"> playwright</span><span style="--shiki-dark:#CE9178;--shiki-light:#A31515"> test</span><span style="--shiki-dark:#569CD6;--shiki-light:#0000FF"> --debug</span></span>
<span class="line"><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000">```</span></span>
<span class="line"></span>
<span class="line"><span style="--shiki-dark:#569CD6;--shiki-dark-font-weight:bold;--shiki-light:#800000;--shiki-light-font-weight:bold">## Best practices</span></span>
<span class="line"></span>
<span class="line"><span style="--shiki-dark:#6796E6;--shiki-light:#0451A5">-</span><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000"> Use data-testid attributes for dynamic content</span></span>
<span class="line"><span style="--shiki-dark:#6796E6;--shiki-light:#0451A5">-</span><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000"> Keep tests independent and atomic</span></span>
<span class="line"><span style="--shiki-dark:#6796E6;--shiki-light:#0451A5">-</span><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000"> Use Page Object Model for complex pages</span></span>
<span class="line"><span style="--shiki-dark:#6796E6;--shiki-light:#0451A5">-</span><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000"> Take screenshots on failure</span></span>
<span class="line"></span></code></pre>
</details>
<details>
<summary>Example: GitHub Actions debugging skill</summary>
<pre class="shiki" data-lang="markdown" shiki-themes dark-plus light-plus" style="--shiki-dark:#D4D4D4;--shiki-light:#000000;--shiki-dark-bg:#1E1E1E;--shiki-light-bg:#FFFFFF" tabindex="0"><code><span class="line"><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000">---</span></span>
<span class="line"><span style="--shiki-dark:#569CD6;--shiki-light:#800000">name</span><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000FF">: </span><span style="--shiki-dark:#CE9178;--shiki-light:#0000FF">github-actions-debugging</span></span>
<span class="line"><span style="--shiki-dark:#569CD6;--shiki-light:#800000">description</span><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000FF">: </span><span style="--shiki-dark:#CE9178;--shiki-light:#0000FF">Guide for debugging failing GitHub Actions workflows. Use this when asked to debug failing GitHub Actions workflows.</span></span>
<span class="line"><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000">---</span></span>
<span class="line"></span>
<span class="line"><span style="--shiki-dark:#569CD6;--shiki-dark-font-weight:bold;--shiki-light:#800000;--shiki-light-font-weight:bold"># GitHub Actions Debugging</span></span>
<span class="line"></span>
<span class="line"><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000">This skill helps you debug failing GitHub Actions workflows in pull requests.</span></span>
<span class="line"></span>
<span class="line"><span style="--shiki-dark:#569CD6;--shiki-dark-font-weight:bold;--shiki-light:#800000;--shiki-light-font-weight:bold">## Process</span></span>
<span class="line"></span>
<span class="line"><span style="--shiki-dark:#6796E6;--shiki-light:#0451A5">1.</span><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000"> Use the </span><span style="--shiki-dark:#CE9178;--shiki-light:#800000">`list_workflow_runs`</span><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000"> tool to look up recent workflow runs for the pull request and their status</span></span>
<span class="line"><span style="--shiki-dark:#6796E6;--shiki-light:#0451A5">2.</span><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000"> Use the </span><span style="--shiki-dark:#CE9178;--shiki-light:#800000">`summarize_job_log_failures`</span><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000"> tool to get an AI summary of the logs for failed jobs</span></span>
<span class="line"><span style="--shiki-dark:#6796E6;--shiki-light:#0451A5">3.</span><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000"> If you need more information, use the </span><span style="--shiki-dark:#CE9178;--shiki-light:#800000">`get_job_logs`</span><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000"> or </span><span style="--shiki-dark:#CE9178;--shiki-light:#800000">`get_workflow_run_logs`</span><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000"> tool to get the full failure logs</span></span>
<span class="line"><span style="--shiki-dark:#6796E6;--shiki-light:#0451A5">4.</span><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000"> Try to reproduce the failure locally in your environment</span></span>
<span class="line"><span style="--shiki-dark:#6796E6;--shiki-light:#0451A5">5.</span><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000"> Fix the failing build and verify the fix before committing changes</span></span>
<span class="line"></span>
<span class="line"><span style="--shiki-dark:#569CD6;--shiki-dark-font-weight:bold;--shiki-light:#800000;--shiki-light-font-weight:bold">## Common issues</span></span>
<span class="line"></span>
<span class="line"><span style="--shiki-dark:#6796E6;--shiki-light:#0451A5">-</span><span style="--shiki-dark:#569CD6;--shiki-dark-font-weight:bold;--shiki-light:#000080;--shiki-light-font-weight:bold"> **Missing environment variables**</span><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000">: Check that all required secrets are configured</span></span>
<span class="line"><span style="--shiki-dark:#6796E6;--shiki-light:#0451A5">-</span><span style="--shiki-dark:#569CD6;--shiki-dark-font-weight:bold;--shiki-light:#000080;--shiki-light-font-weight:bold"> **Version mismatches**</span><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000">: Verify action versions and dependencies are compatible</span></span>
<span class="line"><span style="--shiki-dark:#6796E6;--shiki-light:#0451A5">-</span><span style="--shiki-dark:#569CD6;--shiki-dark-font-weight:bold;--shiki-light:#000080;--shiki-light-font-weight:bold"> **Permission issues**</span><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000">: Ensure the workflow has the necessary permissions</span></span>
<span class="line"><span style="--shiki-dark:#6796E6;--shiki-light:#0451A5">-</span><span style="--shiki-dark:#569CD6;--shiki-dark-font-weight:bold;--shiki-light:#000080;--shiki-light-font-weight:bold"> **Timeout issues**</span><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000">: Consider splitting long-running jobs or increasing timeout values</span></span>
<span class="line"></span></code></pre>
</details>
<h2 id="_use-skills-as-slash-commands" data-needslink="_use-skills-as-slash-commands">Use skills as slash commands</h2>
<p>Skills are available as slash commands in chat, alongside <a href="/docs/copilot/customization/prompt-files">prompt files</a>. Type <code>/</code> in the chat input field to see a list of available skills and prompts, and select a skill to invoke it.</p>
<p>You can add extra context after the slash command. For example, <code>/webapp-testing for the login page</code> or <code>/github-actions-debugging PR #42</code>.</p>
<p>By default, all skills appear in the <code>/</code> menu. Use the <code>user-invocable</code> and <code>disable-model-invocation</code> frontmatter properties to control how each skill is accessed:</p>
<table class="table table-striped">
<thead>
<tr>
<th>Configuration</th>
<th>Slash command</th>
<th>Auto-loaded by Copilot</th>
<th>Use case</th>
</tr>
</thead>
<tbody>
<tr>
<td>Default (both properties omitted)</td>
<td>Yes</td>
<td>Yes</td>
<td>General-purpose skills</td>
</tr>
<tr>
<td><code>user-invocable: false</code></td>
<td>No</td>
<td>Yes</td>
<td>Background knowledge skills that the model loads when relevant</td>
</tr>
<tr>
<td><code>disable-model-invocation: true</code></td>
<td>Yes</td>
<td>No</td>
<td>Skills you only want to run on demand</td>
</tr>
<tr>
<td>Both set</td>
<td>No</td>
<td>No</td>
<td>Disabled skills</td>
</tr>
</tbody>
</table>
<h2 id="_how-copilot-uses-skills" data-needslink="_how-copilot-uses-skills">How Copilot uses skills</h2>
<p>Skills load content progressively to keep your context efficient. Here is an example of how Copilot uses the <code>webapp-testing</code> skill:</p>
<ol>
<li>
<p><strong>Discovery</strong>: Copilot reads the skill's <code>name</code> and <code>description</code> from the YAML frontmatter. When you ask &quot;help me test the login page&quot;, Copilot matches this to the <code>webapp-testing</code> skill based on its description.</p>
</li>
<li>
<p><strong>Instructions loading</strong>: Copilot loads the <code>SKILL.md</code> body into its context, giving it access to the detailed testing procedures and guidelines. You can also trigger this step directly by typing <code>/webapp-testing</code> in chat.</p>
</li>
<li>
<p><strong>Resource access</strong>: As Copilot works through the instructions, it accesses additional files in the skill directory, such as <code>test-template.js</code> or example scenarios, only when it references them. If a file isn't referenced in the instructions, it won't be loaded.</p>
</li>
</ol>
<p>This three-level loading system means you can install many skills without consuming context. Copilot loads only what is relevant for each task.</p>
<h2 id="_use-shared-skills" data-needslink="_use-shared-skills">Use shared skills</h2>
<p>You can use skills created by others to enhance Copilot's capabilities. The <a href="https://github.com/github/awesome-copilot" class="external-link" target="_blank">github/awesome-copilot</a> repository contains a growing community collection of skills, custom agents, instructions, and prompts. The <a href="https://github.com/anthropics/skills" class="external-link" target="_blank">anthropics/skills</a> repository contains additional reference skills.</p>
<p>You can also discover and install skills that are bundled in <a href="/docs/copilot/customization/agent-plugins">agent plugins</a>. Skills from installed plugins appear alongside your locally defined skills in the <strong>Configure Skills</strong> menu.</p>
<p>To use a shared skill:</p>
<ol>
<li>Browse the available skills in the repository</li>
<li>Copy the skill directory to your <code>.github/skills/</code> folder</li>
<li>Review and customize the <code>SKILL.md</code> file for your needs</li>
<li>Optionally, modify or add resources as needed</li>
</ol>
<div class="markdown-alert tip" dir="auto">
      <span>
        <svg class="markdown-alert-icon" viewBox="0 0 16 16" version="1.1" width="16" height="16" aria-hidden="true">
          <path d="M8 1.5c-2.363 0-4 1.69-4 3.75 0 .984.424 1.625.984 2.304l.214.253c.223.264.47.556.673.848.284.411.537.896.621 1.49a.75.75 0 0 1-1.484.211c-.04-.282-.163-.547-.37-.847a8.456 8.456 0 0 0-.542-.68c-.084-.1-.173-.205-.268-.32C3.201 7.75 2.5 6.766 2.5 5.25 2.5 2.31 4.863 0 8 0s5.5 2.31 5.5 5.25c0 1.516-.701 2.5-1.328 3.259-.095.115-.184.22-.268.319-.207.245-.383.453-.541.681-.208.3-.33.565-.37.847a.751.751 0 0 1-1.485-.212c.084-.593.337-1.078.621-1.489.203-.292.45-.584.673-.848.075-.088.147-.173.213-.253.561-.679.985-1.32.985-2.304 0-2.06-1.637-3.75-4-3.75ZM5.75 12h4.5a.75.75 0 0 1 0 1.5h-4.5a.75.75 0 0 1 0-1.5ZM6 15.25a.75.75 0 0 1 .75-.75h2.5a.75.75 0 0 1 0 1.5h-2.5a.75.75 0 0 1-.75-.75Z">
          </path>
        </svg>
        Tip
      </span><p>Always review shared skills before using them to ensure they meet your requirements and security standards. VS Code's <a href="/docs/copilot/agents/agent-tools#_terminal-commands">terminal tool</a> provides controls for script execution, including <a href="/docs/copilot/agents/agent-tools#_automatically-approve-terminal-commands">auto-approve options</a> with configurable allow-lists and tight controls over which code runs. Learn more about <a href="/docs/copilot/security#_automated-approval">security considerations</a> for auto-approval features.</p>
</div><h2 id="_contribute-skills-from-extensions" data-needslink="_contribute-skills-from-extensions">Contribute skills from extensions</h2>
<p>Extensions can contribute skills using the <code>chatSkills</code> contribution point in their <code>package.json</code>. The path must point to a directory that contains a <code>SKILL.md</code> file, following the <a href="https://agentskills.io/specification" class="external-link" target="_blank">Agent Skills specification</a>.</p>
<h3 id="_required-folder-structure" data-needslink="_required-folder-structure">Required folder structure</h3>
<p>The skill directory must follow this structure:</p>
<pre class="shiki" data-lang="text" shiki-themes dark-plus light-plus" style="--shiki-dark:#D4D4D4;--shiki-light:#000000;--shiki-dark-bg:#1E1E1E;--shiki-light-bg:#FFFFFF" tabindex="0"><code><span class="line"><span>extension-root/</span></span>
<span class="line"><span>└── skills/</span></span>
<span class="line"><span>    └── my-skill/           # Directory name must match the `name` field in SKILL.md</span></span>
<span class="line"><span>        └── SKILL.md         # Required</span></span>
<span class="line"><span></span></span></code></pre>
<h3 id="_register-the-skill-in-packagejson" data-needslink="_register-the-skill-in-packagejson">Register the skill in package.json</h3>
<p>Add the <code>chatSkills</code> contribution point in your extension's <code>package.json</code>. The <code>path</code> property must point to the corresponding <code>SKILL.md</code> file:</p>
<pre class="shiki" data-lang="json" shiki-themes dark-plus light-plus" style="--shiki-dark:#D4D4D4;--shiki-light:#000000;--shiki-dark-bg:#1E1E1E;--shiki-light-bg:#FFFFFF" tabindex="0"><code><span class="line"><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000">{</span></span>
<span class="line"><span style="--shiki-dark:#9CDCFE;--shiki-light:#0451A5">  "contributes"</span><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000">: {</span></span>
<span class="line"><span style="--shiki-dark:#9CDCFE;--shiki-light:#0451A5">    "chatSkills"</span><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000">: [</span></span>
<span class="line"><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000">      {</span></span>
<span class="line"><span style="--shiki-dark:#9CDCFE;--shiki-light:#0451A5">        "path"</span><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000">: </span><span style="--shiki-dark:#CE9178;--shiki-light:#A31515">"./skills/my-skill/SKILL.md"</span></span>
<span class="line"><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000">      }</span></span>
<span class="line"><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000">    ]</span></span>
<span class="line"><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000">  }</span></span>
<span class="line"><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000">}</span></span>
<span class="line"></span></code></pre>
<div class="markdown-alert important" dir="auto">
      <span>
        <svg class="markdown-alert-icon" viewBox="0 0 16 16" version="1.1" width="16" height="16" aria-hidden="true">
          <path d="M0 1.75C0 .784.784 0 1.75 0h12.5C15.216 0 16 .784 16 1.75v9.5A1.75 1.75 0 0 1 14.25 13H8.06l-2.573 2.573A1.458 1.458 0 0 1 3 14.543V13H1.75A1.75 1.75 0 0 1 0 11.25Zm1.75-.25a.25.25 0 0 0-.25.25v9.5c0 .138.112.25.25.25h2a.75.75 0 0 1 .75.75v2.19l2.72-2.72a.749.749 0 0 1 .53-.22h6.5a.25.25 0 0 0 .25-.25v-9.5a.25.25 0 0 0-.25-.25Zm7 2.25v2.5a.75.75 0 0 1-1.5 0v-2.5a.75.75 0 0 1 1.5 0ZM9 9a1 1 0 1 1-2 0 1 1 0 0 1 2 0Z">
          </path>
        </svg>
        Important
      </span><p>The <code>name</code> field in the <code>SKILL.md</code> frontmatter must match the parent directory name. For example, if the directory is <code>skills/my-skill/</code>, the <code>name</code> field must be <code>my-skill</code>. If the name does not match, the skill is not loaded.</p>
</div><p>The <code>SKILL.md</code> file follows the same format as <a href="#_create-a-skill">project and personal skills</a>. For example:</p>
<pre class="shiki" data-lang="markdown" shiki-themes dark-plus light-plus" style="--shiki-dark:#D4D4D4;--shiki-light:#000000;--shiki-dark-bg:#1E1E1E;--shiki-light-bg:#FFFFFF" tabindex="0"><code><span class="line"><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000">---</span></span>
<span class="line"><span style="--shiki-dark:#569CD6;--shiki-light:#800000">name</span><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000FF">: </span><span style="--shiki-dark:#CE9178;--shiki-light:#0000FF">my-skill</span></span>
<span class="line"><span style="--shiki-dark:#569CD6;--shiki-light:#800000">description</span><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000FF">: </span><span style="--shiki-dark:#CE9178;--shiki-light:#0000FF">Description of what the skill does and when to use it.</span></span>
<span class="line"><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000">---</span></span>
<span class="line"></span>
<span class="line"><span style="--shiki-dark:#569CD6;--shiki-dark-font-weight:bold;--shiki-light:#800000;--shiki-light-font-weight:bold"># My Skill</span></span>
<span class="line"></span>
<span class="line"><span style="--shiki-dark:#D4D4D4;--shiki-light:#000000">Detailed instructions for the skill...</span></span>
<span class="line"></span></code></pre>
<h2 id="_agent-skills-standard" data-needslink="_agent-skills-standard">Agent Skills standard</h2>
<p>Agent Skills is an open standard that enables portability across different AI agents. Skills you create in VS Code work with multiple agents, including:</p>
<ul>
<li><strong>GitHub Copilot in VS Code</strong>: Available in chat and agent mode</li>
<li><strong>GitHub Copilot CLI</strong>: Accessible when working in the terminal</li>
<li><strong>GitHub Copilot cloud agent</strong>: Used during automated coding tasks</li>
</ul>
<p>Learn more about the Agent Skills standard at <a href="https://agentskills.io" class="external-link" target="_blank">agentskills.io</a>.</p>
<h2 id="_related-resources" data-needslink="_related-resources">Related resources</h2>
<ul>
<li><a href="/docs/copilot/customization/overview">Customize AI responses overview</a></li>
<li><a href="/docs/copilot/customization/custom-instructions">Create custom instructions</a></li>
<li><a href="/docs/copilot/customization/prompt-files">Create reusable prompt files</a></li>
<li><a href="/docs/copilot/customization/custom-agents">Create custom agents</a></li>
<li><a href="https://agentskills.io" class="external-link" target="_blank">Agent Skills specification</a></li>
<li><a href="https://github.com/anthropics/skills" class="external-link" target="_blank">Reference skills repository</a></li>
<li><a href="/docs/copilot/customization/agent-plugins">Discover and manage agent plugins</a></li>
</ul>

                <div class="feedback" data-edit-url="https://vscode.dev/github/microsoft/vscode-docs/blob/main/docs/copilot/customization/agent-skills.md"></div>
                
                <div class="body-footer">4/22/2026</div>
                
            </main>
        </div>
        
        <!-- Mobile connect widget -->
        <div class="docs-mobile-widgets visible-xs">
            <div class="connect-widget"></div>
        </div>
    </div>
</div>
		</main>
	</div>

	<div id="search-popup-overlay" class="search-popup-overlay" role="dialog" aria-modal="true" aria-label="Search">
		<div class="search-popup-container">
			<div class="search-popup-header">
				<form class="search-popup-form">
					<div class="input-group">
						<input type="text" name="q" class="search-box form-control" placeholder="Search the website"
							aria-label="Search text" role="combobox" aria-expanded="false"
							aria-autocomplete="list" aria-controls="search-results-listbox"
							aria-activedescendant="" />
						<span class="input-group-btn">
							<button tabindex="0" class="btn" type="submit" aria-label="Search">
								<img class="search-icon-dark" src="/assets/icons/search-dark.svg" alt="Search" />
								<img class="search-icon-light" src="/assets/icons/search.svg" alt="Search" />
							</button>
						</span>
					</div>
				</form>
				<button class="search-popup-close" type="button" aria-label="Close search"><i class="codicon codicon-close" aria-hidden="true"></i></button>
			</div>
			<div class="search-popup-results">
				<ul class="search-popup-results-list" id="search-results-listbox" role="listbox" aria-label="Search results"></ul>
			</div>
			<div class="sr-only" role="status" aria-live="polite" aria-atomic="true" id="search-popup-status"></div>
		</div>
	</div>


	<footer role="contentinfo" class="container">
		<div class="footer-container">
			<div class="footer-row">
				<div class="footer-social">
					<ul class="links">
						<li>
							<a href="https://github.com/microsoft/vscode"><img src="/assets/icons/github-icon.svg" alt="VS Code on Github"></a>
						</li>
						<li>
							<a href="https://go.microsoft.com/fwlink/?LinkID=533687"><img src="/assets/icons/x-icon.svg" class="x-icon" alt="Follow us on X"></a>
						</li>
						<li>
							<a href="https://www.linkedin.com/showcase/vs-code"><img src="/assets/icons/linkedin-icon.svg" alt="VS Code on LinkedIn"></a>
						</li>
						<li>
							<a href="https://bsky.app/profile/vscode.dev"><img src="/assets/icons/bluesky-icon.svg" alt="VS Code on Bluesky"></a>
						</li>
						<li>
							<a href="https://www.reddit.com/r/vscode/"><img src="/assets/icons/reddit-icon.svg" alt="Join the VS Code community on Reddit"></a>
						</li>
						<li>
							<a href="https://www.vscodepodcast.com"><img src="/assets/icons/podcast-icon.svg" alt="The VS Code Insiders Podcast"></a>
						</li>
						<li>
							<a href="https://www.tiktok.com/@vscode"><img src="/assets/icons/tiktok-icon.svg" alt="VS Code on TikTok"></a>
						</li>
						<li>
							<a href="https://www.youtube.com/@code"><img src="/assets/icons/youtube-icon.svg" alt="VS Code on YouTube"></a>
						</li>
						<script>
							function manageConsent() {
								if (siteConsent && siteConsent.isConsentRequired) {
									siteConsent.manageConsent();
								}
							}
						</script>
					</ul>
					<a id="footer-microsoft-link" class="microsoft-logo" href="https://www.microsoft.com">
						<img src="/assets/icons/microsoft.svg" alt="Microsoft homepage" />
					</a>
				</div>
			</div>
			<div class="footer-row">
				<ul class="links">
					<li><a id="footer-support-link" href="https://support.serviceshub.microsoft.com/supportforbusiness/create?sapId=d66407ed-3967-b000-4cfb-2c318cad363d"
						target="_blank" rel="noopener" title="Get support for VS Code"
						aria-label="Get support for VS Code (opens in new tab)">Support</a></li>
					<li><a id="footer-privacy-link" href="https://go.microsoft.com/fwlink/?LinkId=521839"
						target="_blank" rel="noopener" title="View the Microsoft privacy statement"
						aria-label="Microsoft privacy statement (opens in new tab)">Privacy</a></li>
					<li style="display: none;"><a id="footer-cookie-link" style="cursor: pointer;" onclick="manageConsent()"
						target="_blank" rel="noopener">Manage Cookies</a></li>
					<li><a id="footer-terms-link" href="https://www.microsoft.com/legal/terms-of-use"
						target="_blank" rel="noopener" title="View the Microsoft Terms of Use"
						aria-label="Microsoft Terms of Use (opens in new tab)">Terms of Use</a></li>
					<li><a id="footer-license-link" href="/License"
						target="_blank" rel="noopener" title="View the Visual Studio Code license"
						aria-label="Visual Studio Code license (opens in new tab)">License</a></li>
				</ul>
			</div>
			<div class="footer-row">
				<ul class="links">
					<li>
						<svg class="privacy-choices" role="img" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 30 14" xml:space="preserve" height="16" width="43">
							<title>Your Privacy Choices Opt-Out Icon</title>
							<path d="M7.4 12.8h6.8l3.1-11.6H7.4C4.2 1.2 1.6 3.8 1.6 7s2.6 5.8 5.8 5.8z" style="fill-rule:evenodd;clip-rule:evenodd;fill:#fff"></path>
							<path d="M22.6 0H7.4c-3.9 0-7 3.1-7 7s3.1 7 7 7h15.2c3.9 0 7-3.1 7-7s-3.2-7-7-7zm-21 7c0-3.2 2.6-5.8 5.8-5.8h9.9l-3.1 11.6H7.4c-3.2 0-5.8-2.6-5.8-5.8z" style="fill-rule:evenodd;clip-rule:evenodd;fill:#06f"></path>
							<path d="M24.6 4c.2.2.2.6 0 .8L22.5 7l2.2 2.2c.2.2.2.6 0 .8-.2.2-.6.2-.8 0l-2.2-2.2-2.2 2.2c-.2.2-.6.2-.8 0-.2-.2-.2-.6 0-.8L20.8 7l-2.2-2.2c-.2-.2-.2-.6 0-.8.2-.2.6-.2.8 0l2.2 2.2L23.8 4c.2-.2.6-.2.8 0z" style="fill:#fff"></path>
							<path d="M12.7 4.1c.2.2.3.6.1.8L8.6 9.8c-.1.1-.2.2-.3.2-.2.1-.5.1-.7-.1L5.4 7.7c-.2-.2-.2-.6 0-.8.2-.2.6-.2.8 0L8 8.6l3.8-4.5c.2-.2.6-.2.9 0z" style="fill:#06f"></path>
						</svg>
						<a id="footer-privacy-choices-link" href="https://aka.ms/YourCaliforniaPrivacyChoices"
						target="_blank" rel="noopener" title="View Your Privacy Choices"
						aria-label="Your Privacy Choices (opens in new tab)">Your Privacy Choices</a></li>
					<li><a id="footer-consumer-health-privacy-link" href="https://go.microsoft.com/fwlink/?linkid=2259814"
						target="_blank" rel="noopener" title="View the Microsoft Consumer Health Privacy policy"
						aria-label="Microsoft Consumer Health Privacy policy (opens in new tab)">Consumer Health Privacy</a></li>
				</ul>
			</div>
		</div>
	</footer>
	<script type="module">
		document.addEventListener('DOMContentLoaded', () => {
			const copilotDeepLinks = document.querySelectorAll('.copilot-deep-link');
			if (copilotDeepLinks.length === 0) {
				return;
			}
			if (window.innerWidth < 992) {
				for (const link of copilotDeepLinks) {
					link.href = 'https://aka.ms/vscode-activatecopilotfree';
				}
			}
		});
	</script>

	<script src="/dist/index.js"></script>

	

	<script type="application/ld+json">
		{
			"@context" : "http://schema.org",
			"@type" : "SoftwareApplication",
			"name" : "Visual Studio Code",
			"softwareVersion": "1.117",
			"offers": {
				"@type": "Offer",
				"price": "0",
				"priceCurrency": "USD"
			},
			"applicationCategory": "DeveloperApplication",
			"applicationSubCategory": "Text Editor",
			"alternateName": "VS Code",
			"datePublished": "2021-11-03",
			"operatingSystem": "Mac, Linux, Windows",
			"logo": "https://code.visualstudio.com/assets/apple-touch-icon.png",
			"screenshot": "https://code.visualstudio.com/assets/images/product-screenshot.png",
			"releaseNotes": "https://code.visualstudio.com/updates",
			"downloadUrl": "https://code.visualstudio.com/download",
			"license": "https://code.visualstudio.com/license",
			"softwareRequirements": "https://code.visualstudio.com/docs/supporting/requirements",
			"url" : "https://code.visualstudio.com",
			"author": {
				"@type": "Organization",
				"name": "Microsoft"
			},
			"publisher": {
				"@type": "Organization",
				"name": "Microsoft"
			},
			"maintainer": {
				"@type": "Organization",
				"name": "Microsoft"
			},
			"potentialAction": {
				"@type": "SearchAction",
				"target": "https://code.visualstudio.com/Search?q={search_term_string}",
				"query-input": "required name=search_term_string"
			},
			"sameAs" : [
				"https://en.wikipedia.org/wiki/Visual_Studio_Code",
				"https://twitter.com/code",
				"https://www.youtube.com/code",
				"https://www.tiktok.com/@vscode",
				"https://github.com/microsoft/vscode"
			]
		}
	</script>
</body>

</html>