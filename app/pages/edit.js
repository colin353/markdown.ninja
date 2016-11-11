/*
  edit.js
  @flow

  The site editor page.
*/

var React = require('react');
import { ContextMenu, MenuItem, ContextMenuTrigger } from "react-contextmenu";

import type { APIInstance } from '../api/api';
declare var api: APIInstance;

var Tree = require('../components/tree');
var Button = require('../components/button');
var Editor = require('../components/editor');
var Preview = require('../components/preview');
var Tab = require('../components/tab');
var Popover = require('../components/popover');

class Edit extends React.Component {
  state: {
    markdown: string,
    html: string,
    pages: Array<Page>,
    showEditor: boolean,
    showPreview: boolean,
    unsavedChanges: boolean,
    selectedPage: Page,
    showRenamePopover: boolean,
    renameValue: string,
    contextPage: Page,
    showDeletePopover: boolean
  };
  renameInput: any;
  converter: any;

  constructor(props: any) {
    super(props);

    this.state = {
      html: "",
      markdown: "",
      pages: [],
      showEditor: true,
      showPreview: true,
      selectedPage: {
        name: "index.md",
        markdown: "",
        html: ""
      },
      unsavedChanges: false,
      showRenamePopover: false,
      contextPage: {"name": "index.md", "markdown": "", "html": ""},
      renameValue: "index.md",
      showDeletePopover: false
    }

    this.converter = new window.showdown.Converter();
  }

  componentDidMount() {
    this.getPages().then((pages) => {
      this.setState({
        selectedPage: pages[0],
        pages
      });
    });

    api.addListener("ctrl+s", "editor", () => {
      this.save();
    });
  }

  getPages() {
    // Load the list of pages for this site.
    return api.pages().then((pages) => {
      // Two actions: sort the list alphabetically, and then make sure that
      // index.md shows up first.
      return pages.sort((a, b) => {
        if(a.name == "index.md") return -1;
        if(b.name == "index.md") return 1;

        return b.name.localeCompare(a.name);
      });
    })
  }

  save() {
    this.state.selectedPage.markdown = this.state.markdown;
    this.state.selectedPage.html = this.state.html;
    return api.editPage(this.state.selectedPage).then(() => {
      this.setState({unsavedChanges: false});
    });
  }

  textEdited(markdown: string) {
    var html = this.converter.makeHtml(markdown);
    this.setState({
      markdown,
      html,
      unsavedChanges:(this.state.selectedPage.markdown != markdown)
    });
  }

  clickTab(tab: 'editor' | 'preview') {
    if(tab == 'editor') {
      // Shouldn't be able to hide the editor if there isn't
      // anything else being displayed.
      if(this.state.showPreview || !this.state.showEditor)
        this.setState({showEditor: !this.state.showEditor})
    } else if (tab == 'preview') {
      if(!this.state.showPreview || this.state.showEditor)
        this.setState({showPreview: !this.state.showPreview})
    }
  }

  clickPage(page: Page) {
    // Don't do anything if we are already displaying that page.
    if(page.name == this.state.selectedPage) return;

    // Save the existing page, then load the new one.
    this.save().then(() => {
      return api.getPage(page.name);
    }).then((page) => {
      this.setState({
        markdown: page.markdown,
        html: page.html,
        selectedPage: page
      });
    })
  }
  handleClick(button: string, e: any, data: { page: Page }) {
    if(button == "rename") {
      this.setState({
        showRenamePopover: true,
        contextPage: data.page,
        renameValue: data.page.name
      });
    }
    else if(button == "delete") {
      this.setState({
        showDeletePopover: true,
        contextPage: data.page
      });
    }
  }

  convertToValidFilename(s: string) {
    var stripped = s.replace(/[^A-Za-z0-9_\\.]+/g, '');
    return stripped;
  }

  renameKeyPress(e: any) {
    // If they hit the enter key, we'll change the name.
    var newName = this.state.renameValue;
    if(e.key == "Enter") {
      api.renamePage(this.state.contextPage.name, newName).then(() => {
        this.setState({showRenamePopover: false});
        return this.getPages();
      }).then((pages) => {
        // If we are currently editing the file that was renamed, let's
        // make sure we update the name there too.
        if(this.state.selectedPage.name == this.state.contextPage.name) {
          this.state.selectedPage.name = newName;
        }
        this.setState({pages, selectedPage: this.state.selectedPage});
      });
    }
  }

  typeRename(e: any) {
    this.setState({
      renameValue: this.convertToValidFilename(e.target.value)
    });
  }

  deletePage(page: Page) {
    api.deletePage(page).then(() => {
      this.setState({showDeletePopover: false});
      return this.getPages();
    }).then((pages) => {
      // If we are currently editing the file that was deleted, let's
      // pull up another one instead.
      if(this.state.selectedPage.name == page.name) {
        this.setState({pages, selectedPage: pages[0]});
      } else {
        this.setState({pages});
      }
    });
  }

  addNewPage() {
    var page: Page;
    api.createPage({markdown: "", html: ""}).then((p) => {
      page = p;
      return this.getPages();
    }).then((pages) => {
      this.setState({pages, selectedPage: page});
    });
  }

  uploadFile() {

  }

  render() {
    return (
      <div style={styles.container}>
        <Tree
          onAddNewPage={this.addNewPage.bind(this)}
          onUploadFile={this.uploadFile.bind(this)}
          clickPage={this.clickPage.bind(this)}
          pages={this.state.pages}
        />
        <div style={styles.editWindow}>
          <div style={styles.tabs}>
            <Tab onClick={this.clickTab.bind(this, 'editor')} selected={this.state.showEditor} indicator={this.state.unsavedChanges} name={this.state.selectedPage.name} />
            <Tab onClick={this.clickTab.bind(this, 'preview')} selected={this.state.showPreview} name="preview" />
          </div>
          <div style={styles.window}>
            <Editor visible={this.state.showEditor} initialText={this.state.selectedPage.markdown} onChange={this.textEdited.bind(this)} />
            {this.state.showPreview?(
              <Preview html={this.state.html} />
            ):[]}
          </div>
        </div>
        <ContextMenu id="page">
          <MenuItem key="rename" onClick={this.handleClick.bind(this, "rename")}>
            Rename
          </MenuItem>
          <MenuItem key="delete" onClick={this.handleClick.bind(this, "delete")}>
            Delete
          </MenuItem>
          <MenuItem divider />
          <MenuItem key="newfile" onClick={this.handleClick.bind(this, "newfile")}>
            New file
          </MenuItem>
          <MenuItem key="upload" onClick={this.handleClick.bind(this, "upload")}>
            Upload file
          </MenuItem>
        </ContextMenu>

        <Popover
          onFocus={() => { this.renameInput.focus(); this.renameInput.setSelectionRange(0, this.state.renameValue.length-3); }}
          visible={this.state.showRenamePopover}
          onDismiss={() => this.setState({showRenamePopover: false})}
        >
          <p>Renaming page from <b>{this.state.contextPage.name}</b> to: </p>
          <input
            ref={(r) => this.renameInput = r}
            style={{width: 380}}
            type="text"
            value={this.state.renameValue}
            onKeyPress={this.renameKeyPress.bind(this)}
            onChange={this.typeRename.bind(this)}
          />
        </Popover>

        <Popover
          visible={this.state.showDeletePopover}
          onDismiss={() => this.setState({showDeletePopover: false})}
        >
          <p>So you want to delete <b>{this.state.contextPage.name}</b>?</p>

          <div style={{display: 'flex'}}>
            <Button action="yes" color="red" onClick={this.deletePage.bind(this,this.state.contextPage)} />
            <div style={{marginLeft: 20}}></div>
            <Button action="no" onClick={() => this.setState({showDeletePopover: false})} />
          </div>
        </Popover>
      </div>
    );
  }
}

const styles = {
  container: {
    display: 'flex',
    flexDirection: 'row',
    flex: 1
  },
  editWindow: {
    display: 'flex',
    flexDirection: 'column',
    flex: 1
  },
  tabs: {
    display: 'flex',
    flexDirection: 'row',
    backgroundColor: '#716669'
  },
  window: {
    flex: 1,
    display: 'flex',
    flexDirection: 'row'
  }
};

module.exports = Edit;
