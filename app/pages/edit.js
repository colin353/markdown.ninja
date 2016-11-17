/*
  edit.js
  @flow

  The site editor page.
*/

var React = require('react');
import { ContextMenu, MenuItem, ContextMenuTrigger } from "react-contextmenu";

var Showdown = require('showdown');
var Tree = require('../components/tree');
var Button = require('../components/button');
var Editor = require('../components/editor');
var Preview = require('../components/preview');
var Tab = require('../components/tab');
var Popover = require('../components/popover');
var ConvertAbsToRelativeURL = require('../tools/relative-to-absolute');

import type { APIInstance } from '../api/api';

type Props = {
  api: APIInstance
};

class Edit extends React.Component {
  props: Props;
  state: {
    markdown: string,
    html: string,
    pages: Array<Page>,
    files: Array<File>,
    showEditor: boolean,
    showPreview: boolean,
    unsavedChanges: boolean,
    selectedPage: Page,
    showRenamePopover: boolean,
    renameType: string,
    renameValue: string,
    contextPage: Page,
    contextFile: File,
    showDeletePagePopover: boolean,
    showDeleteFilePopover: boolean,
    showUploadPopover: boolean
  };
  renameInput: any;
  fileInput: any;
  converter: any;

  constructor(props: any) {
    super(props);

    this.state = {
      html: "",
      markdown: "",
      pages: [],
      files: [],
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
      contextFile: {"name": ""},
      renameValue: "index.md",
      showDeletePagePopover: false,
      showDeleteFilePopover: false,
      showUploadPopover: false,
      renameType: 'page'
    }

    this.converter = new Showdown.Converter();
  }

  componentDidMount() {
    this.getPages().then((pages) => {
      this.setState({
        selectedPage: pages[0],
        pages
      });
    });

    this.getFiles().then((files) => {
      this.setState({files});
    })

    this.props.api.addListener("ctrl+s", "editor", () => {
      this.save();
    });
  }

  getFiles() {
    return this.props.api.files().then((files) => {
      // Sort the list alphabetically.
      return files.sort((a, b) => {
        return b.name.localeCompare(a.name);
      })
    });
  }

  getPages() {
    // Load the list of pages for this site.
    return this.props.api.pages().then((pages) => {
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
    return this.props.api.editPage(this.state.selectedPage).then(() => {
      this.setState({unsavedChanges: false});
    });
  }

  textEdited(markdown: string) {
    var html = this.converter.makeHtml(markdown);
    // We also need to replace relative URLs with absolute URLs, so that
    // the preview window works correctly.
    html = ConvertAbsToRelativeURL(html);

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

  clickFile(file: File) {

  }

  clickPage(page: Page) {
    // Don't do anything if we are already displaying that page.
    if(page.name == this.state.selectedPage.name) return;

    // Save the existing page, then load the new one.
    this.save().then(() => {
      return this.props.api.getPage(page.name);
    }).then((page) => {
      this.setState({
        markdown: page.markdown,
        html: page.html,
        selectedPage: page
      });
    })
  }

  handleClick(button: string, e: any, data: { page?: Page, file?: File }) {
    if(button == "rename" && data.page) {
      this.setState({
        renameType: 'page',
        showRenamePopover: true,
        contextPage: data.page,
        renameValue: data.page.name
      });
    } else if(button == "rename" && data.file) {
      this.setState({
        renameType: 'file',
        showRenamePopover: true,
        contextFile: data.file,
        renameValue: data.file.name
      })
    } else if(button == "delete" && data.file) {
      this.setState({
        showDeleteFilePopover: true,
        contextFile: data.file
      });
    } else if(button == "delete" && data.page) {
      this.setState({
        showDeletePagePopover: true,
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
    if(e.key == "Enter" && this.state.renameType == 'page') {
      this.props.api.renamePage(this.state.contextPage.name, newName).then(() => {
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
    } else if(e.key == "Enter" && this.state.renameType == 'file') {
      this.props.api.renameFile(this.state.contextFile.name, newName).then(() => {
        this.setState({showRenamePopover: false});
        return this.getFiles();
      }).then((files) => {
        this.setState({files});
      });
    }
  }

  typeRename(e: any) {
    this.setState({
      renameValue: this.convertToValidFilename(e.target.value)
    });
  }

  deletePage(page: Page) {
    this.props.api.deletePage(page).then(() => {
      this.setState({showDeletePagePopover: false});
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

  deleteFile(file: File) {
    this.props.api.deleteFile(file.name).then(() => {
      return this.getFiles();
    }).then((files) => {
      this.setState({files, showDeleteFilePopover: false});
    })
  }

  addNewPage() {
    var page: Page;
    this.props.api.createPage({markdown: "", html: ""}).then((p) => {
      page = p;
      return this.getPages();
    }).then((pages) => {
      this.setState({pages, selectedPage: page});
    });
  }

  uploadFile() {
    this.props.api.uploadFile("file.txt", this.fileInput.files[0], () => {}).then(() => {
      return this.getFiles();
    }).then((files) => {
      this.setState({files, showUploadPopover: false});
    });
  }

  clickUploadFile() {
    this.setState({showUploadPopover: true});
  }

  render() {
    return (
      <div style={styles.container}>
        <Tree
          api={this.props.api}
          onAddNewPage={this.addNewPage.bind(this)}
          onUploadFile={this.clickUploadFile.bind(this)}
          clickPage={this.clickPage.bind(this)}
          clickFile={this.clickFile.bind(this)}
          pages={this.state.pages}
          files={this.state.files}
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
          api={this.props.api}
          onFocus={() => { this.renameInput.focus(); this.renameInput.setSelectionRange(0, this.state.renameValue.length-3); }}
          visible={this.state.showRenamePopover}
          onDismiss={() => this.setState({showRenamePopover: false})}
        >
          <p>Renaming {this.state.renameType} from <b>{this.state.renameType=='page'?this.state.contextPage.name:this.state.contextFile.name}</b> to: </p>
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
          api={this.props.api}
          visible={this.state.showDeletePagePopover||this.state.showDeleteFilePopover}
          onDismiss={() => this.setState({showDeleteFilePopover: false, showDeletePagePopover: false})}
        >
          <p>So you want to delete <b>{this.state.showDeletePagePopover?this.state.contextPage.name:this.state.contextFile.name}</b>?</p>

          <div style={{display: 'flex'}}>
            <Button action="yes" color="red" onClick={this.state.showDeletePagePopover?this.deletePage.bind(this,this.state.contextPage):this.deleteFile.bind(this, this.state.contextFile)} />
            <div style={{marginLeft: 20}}></div>
            <Button action="no" onClick={() => this.setState({showDeletePagePopover: false, showDeleteFilePopover: false})} />
          </div>
        </Popover>

        <Popover
          api={this.props.api}
          visible={this.state.showUploadPopover}
          onDismiss={() => this.setState({showUploadPopover: false})}
        >
          <p>Upload a file?</p>
          <div style={{display: 'flex', alignItems: 'center'}}>
            <input ref={(r) => this.fileInput = r} style={{flex: 1}} type="file" />
            <Button onClick={this.uploadFile.bind(this)} color="red" action="upload" />
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
